package service

import (
	"fmt"
	"time"
	"user-center/internal/model"
	"user-center/internal/repository"

	"github.com/xuri/excelize/v2"
)

type AttendanceService struct {
	repo   *repository.AttendanceRepository
	ccRepo *repository.CCRepository
}

func NewAttendanceService() *AttendanceService {
	return &AttendanceService{
		repo:   repository.NewAttendanceRepository(),
		ccRepo: repository.NewCCRepository(),
	}
}

// AttendanceInfo 在班信息
type AttendanceInfo struct {
	CCID             int64  `json:"ccId"`
	CCName           string `json:"ccName"`
	NickName         string `json:"nickName"`
	SquadName        string `json:"squadName"`
	TeamName         string `json:"teamName"`
	LegionName       string `json:"legionName"`
	AttendanceStatus string `json:"attendanceStatus"`
}

// GetAttendanceList 获取在班列表
func (s *AttendanceService) GetAttendanceList(query *model.CCQuery, dates []string) ([]map[string]interface{}, error) {
	// 获取CC列表
	result, err := s.ccRepo.List(query)
	if err != nil {
		return nil, err
	}

	ccList := result.Rows.([]model.CCMember)
	var attendanceList []map[string]interface{}

	for _, cc := range ccList {
		info := map[string]interface{}{
			"ccId":             cc.ID,
			"ccName":           cc.Name,
			"nickName":         cc.NickName,
			"squadName":        cc.SquadName,
			"teamName":         cc.TeamName,
			"legionName":       cc.LegionName,
			"attendanceStatus": cc.AttendanceStatus,
		}

		// 获取各日期的在班状态
		for _, dateStr := range dates {
			date, _ := time.Parse("2006-01-02", dateStr)
			record, _ := s.repo.GetByCCIDAndDate(cc.ID, date)
			if record != nil {
				info[dateStr] = record.Status
			} else {
				info[dateStr] = "" // 未设置
			}
		}

		attendanceList = append(attendanceList, info)
	}

	return attendanceList, nil
}

// UpdateAttendance 更新在班状态
func (s *AttendanceService) UpdateAttendance(ccID int64, dateStr string, status string, operatorID int64) error {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	record := &model.CCAttendance{
		CCID:           ccID,
		AttendanceDate: date,
		Status:         status,
		OperatorID:     &operatorID,
	}

	if err := s.repo.CreateOrUpdate(record); err != nil {
		return err
	}

	// 同时更新CC成员的当前在班状态（如果是今天的记录）
	today := time.Now().Format("2006-01-02")
	if dateStr == today {
		s.ccRepo.UpdateAttendanceStatus(ccID, status)
	}

	return nil
}

// BatchUpdateAttendance 批量更新在班状态
func (s *AttendanceService) BatchUpdateAttendance(updates []model.AttendanceBatchUpdateItem, operatorID int64) error {
	for _, item := range updates {
		if err := s.UpdateAttendance(item.CCID, item.Date, item.Status, operatorID); err != nil {
			return err
		}
	}
	return nil
}

// GetAttendanceStats 获取在班统计
func (s *AttendanceService) GetAttendanceStats(dateStr string) (map[string]int64, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return s.repo.CountByStatus(date)
}

// GetCCAttendanceHistory 获取CC的在班历史
func (s *AttendanceService) GetCCAttendanceHistory(ccID int64, startDate, endDate string) ([]*model.CCAttendance, error) {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	return s.repo.GetByCCID(ccID, start, end)
}

// ExportAttendance 导出在班记录为Excel
func (s *AttendanceService) ExportAttendance(query *model.CCQuery, dates []string) ([]byte, error) {
	// 获取在班列表
	list, err := s.GetAttendanceList(query, dates)
	if err != nil {
		return nil, err
	}

	// 创建Excel文件
	f := excelize.NewFile()
	sheetName := "在班记录"
	f.SetSheetName("Sheet1", sheetName)

	// 设置表头
	headers := []string{"姓名", "昵称", "战队", "团队", "军团"}
	for _, date := range dates {
		headers = append(headers, date)
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
	}

	// 填充数据
	statusMap := map[string]string{
		"1": "在",
		"2": "休",
		"3": "假",
		"":  "--",
	}

	for row, item := range list {
		rowNum := row + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), item["ccName"])
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), item["nickName"])
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), item["squadName"])
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), item["teamName"])
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), item["legionName"])

		for i, date := range dates {
			cell := fmt.Sprintf("%s%d", string(rune('F'+i)), rowNum)
			status := ""
			if v, ok := item[date]; ok && v != nil {
				status = fmt.Sprintf("%v", v)
			}
			f.SetCellValue(sheetName, cell, statusMap[status])
		}
	}

	// 返回文件字节
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
