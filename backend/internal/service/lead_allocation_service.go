package service

import (
	"time"
	"user-center/internal/model"
	"user-center/internal/repository"

	"gorm.io/gorm"
)

type LeadAllocationService struct {
	repo           *repository.LeadAllocationRepository
	ccRepo         *repository.CCRepository
	attendanceRepo *repository.AttendanceRepository
}

func NewLeadAllocationService() *LeadAllocationService {
	return &LeadAllocationService{
		repo:           repository.NewLeadAllocationRepository(),
		ccRepo:         repository.NewCCRepository(),
		attendanceRepo: repository.NewAttendanceRepository(),
	}
}

// LeadAllocationInfo 例子分配信息
type LeadAllocationInfo struct {
	CCID               int64  `json:"ccId"`
	CCName             string `json:"ccName"`
	NickName           string `json:"nickName"`
	SquadName          string `json:"squadName"`
	TeamName           string `json:"teamName"`
	LegionName         string `json:"legionName"`
	AttendanceStatus   string `json:"attendanceStatus"`
	ExpectedAllocation int    `json:"expectedAllocation"`
	ActualAllocation   int    `json:"actualAllocation"`
	ExpectedSupplement int    `json:"expectedSupplement"`
	ActualSupplement   int    `json:"actualSupplement"`
	Overdraft          int    `json:"overdraft"`
	ProcessedOverdraft int    `json:"processedOverdraft"`
	PendingOverdraft   int    `json:"pendingOverdraft"`
	IsAllocated        string `json:"isAllocated"`
	AllocationRule     string `json:"allocationRule"`
	AllocationReason   string `json:"allocationReason"`
}

// GetAllocationList 获取分配列表
func (s *LeadAllocationService) GetAllocationList(query *model.LeadAllocationQuery, dateStr string) (*model.PageResult, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	// 构建CC查询
	ccQuery := &model.CCQuery{
		PageQuery:        query.PageQuery,
		CCID:             query.CCID,
		Name:             query.CCName,
		SquadID:          query.SquadID,
		TeamID:           query.TeamID,
		LegionID:         query.LegionID,
		AttendanceStatus: query.AttendanceStatus,
	}

	// 获取CC列表
	ccResult, err := s.ccRepo.List(ccQuery)
	if err != nil {
		return nil, err
	}

	ccList := ccResult.Rows.([]model.CCMember)
	var allocations []LeadAllocationInfo

	for _, cc := range ccList {
		allocation, _ := s.repo.GetByCCIDAndDate(cc.ID, date)

		info := LeadAllocationInfo{
			CCID:             cc.ID,
			CCName:           cc.Name,
			NickName:         cc.NickName,
			SquadName:        cc.SquadName,
			TeamName:         cc.TeamName,
			LegionName:       cc.LegionName,
			AttendanceStatus: cc.AttendanceStatus,
		}

		if allocation != nil {
			info.ExpectedAllocation = allocation.ExpectedAllocation
			info.ActualAllocation = allocation.ActualAllocation
			info.ExpectedSupplement = allocation.ExpectedSupplement
			info.ActualSupplement = allocation.ActualSupplement
			info.Overdraft = allocation.Overdraft
			info.ProcessedOverdraft = allocation.ProcessedOverdraft
			info.PendingOverdraft = allocation.PendingOverdraft
			info.IsAllocated = allocation.IsAllocated
			info.AllocationRule = allocation.AllocationRule
			info.AllocationReason = allocation.AllocationReason
		}

		allocations = append(allocations, info)
	}

	return &model.PageResult{
		Rows:  allocations,
		Total: ccResult.Total,
	}, nil
}

// UpdateAllocation 更新分配记录
func (s *LeadAllocationService) UpdateAllocation(ccID int64, dateStr string, dto *model.LeadAllocationUpdateDTO) error {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	record := &model.CCLeadAllocation{
		CCID:               ccID,
		AllocationDate:     date,
		ExpectedAllocation: dto.ExpectedAllocation,
		ActualAllocation:   dto.ActualAllocation,
		ExpectedSupplement: dto.ExpectedSupplement,
		ActualSupplement:   dto.ActualSupplement,
		Overdraft:          dto.Overdraft,
		ProcessedOverdraft: dto.ProcessedOverdraft,
		PendingOverdraft:   dto.PendingOverdraft,
		IsAllocated:        dto.IsAllocated,
		AllocationRule:     dto.AllocationRule,
		AllocationReason:   dto.AllocationReason,
	}

	return s.repo.CreateOrUpdate(record)
}

// BatchUpdateIsAllocated 批量更新是否分配状态
func (s *LeadAllocationService) BatchUpdateIsAllocated(ccIDs []int64, dateStr string, isAllocated string, reason string) error {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	for _, ccID := range ccIDs {
		existing, _ := s.repo.GetByCCIDAndDate(ccID, date)
		if existing == nil {
			existing = &model.CCLeadAllocation{
				CCID:           ccID,
				AllocationDate: date,
			}
		}
		existing.IsAllocated = isAllocated
		existing.AllocationReason = reason
		if err := s.repo.CreateOrUpdate(existing); err != nil {
			return err
		}
	}
	return nil
}

// GetAllocationStats 获取分配统计
func (s *LeadAllocationService) GetAllocationStats(dateStr string) (map[string]int64, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return s.repo.SumByDate(date)
}

// AllocationDetailInfo 分配详情信息
type AllocationDetailInfo struct {
	AllocationDate     string `json:"allocationDate"`
	ExpectedAllocation int    `json:"expectedAllocation"`
	ActualAllocation   int    `json:"actualAllocation"`
	Overdraft          int    `json:"overdraft"`
	IsAllocated        string `json:"isAllocated"`
	AllocationRule     string `json:"allocationRule"`
	AllocationReason   string `json:"allocationReason"`
}

// GetAllocationDetail 获取CC分配详情历史
func (s *LeadAllocationService) GetAllocationDetail(ccID int64, startDate, endDate string) ([]AllocationDetailInfo, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	records, err := s.repo.GetByCCID(ccID, start, end)
	if err != nil {
		return nil, err
	}

	var details []AllocationDetailInfo
	for _, r := range records {
		details = append(details, AllocationDetailInfo{
			AllocationDate:     r.AllocationDate.Format("2006-01-02"),
			ExpectedAllocation: r.ExpectedAllocation,
			ActualAllocation:   r.ActualAllocation,
			Overdraft:          r.Overdraft,
			IsAllocated:        r.IsAllocated,
			AllocationRule:     r.AllocationRule,
			AllocationReason:   r.AllocationReason,
		})
	}
	return details, nil
}

// GetSingleDetail 获取分配详情单日明细
func (s *LeadAllocationService) GetSingleDetail(ccID int64, dateStr string) (*model.LeadAllocationSingleDetailDTO, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	// 1. 获取CC信息
	ccMember, err := s.ccRepo.Get(ccID)
	if err != nil {
		return nil, err
	}

	// 2. 获取分配信息
	allocation, err := s.repo.GetByCCIDAndDate(ccID, date)
	if err != nil {
		return nil, err
	}

	// 3. 查找上一个工作日
	lastWorkDayRecord, err := s.attendanceRepo.GetLastWorkDay(ccID, date)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 4. 构建DTO
	dto := &model.LeadAllocationSingleDetailDTO{
		CCID:               ccMember.ID,
		CCName:             ccMember.Name,
		NickName:           ccMember.NickName,
		Date:               dateStr,
		DayOfWeek:          weekdayMap[date.Weekday()],
		AllocationTime:     "14:00", // Mock
		IsAllocated:        "0",
		AttendanceStatus:   ccMember.AttendanceStatus,
		CallDurationTarget: ccMember.CallDurationTarget,
		CallCountTarget:    ccMember.CallCountTarget,
	}

	if allocation != nil {
		dto.AllocationRule = allocation.AllocationRule
		dto.AllocationReason = allocation.AllocationReason
		dto.IsAllocated = allocation.IsAllocated
		dto.ExpectedAllocation = allocation.ExpectedAllocation
		dto.ActualAllocation = allocation.ActualAllocation
		dto.Overdraft = allocation.Overdraft
	}

	if lastWorkDayRecord != nil {
		dto.LastWorkDay = lastWorkDayRecord.AttendanceDate.Format("2006-01-02")
		// Mock 上一工作日数据
		dto.LastWorkDayCallCnt = 21   // Mock
		dto.LastWorkDayCallDur = 7680 // Mock 2h 8m
		dto.LastWorkDayReach = true   // Mock
	}

	// 5. 判断是否为SuperCC (军团长/团长) 并添加 breakdown
	// RoleType: legion_leader (军团长), team_leader (团长), squad_leader (战队长), cc (Cc)
	if ccMember.RoleType == model.RoleTypeLegionLeader || ccMember.RoleType == model.RoleTypeTeamLeader {
		dto.LevelBreakdown = []model.Level{
			{Name: "A", Predicted: 5, Overdraft: 1, Actual: 5},
			{Name: "B", Predicted: 5, Overdraft: 1, Actual: 5},
			{Name: "C", Predicted: 5, Overdraft: 1, Actual: 5},
		}
	}

	return dto, nil
}

var weekdayMap = map[time.Weekday]string{
	time.Sunday:    "周日",
	time.Monday:    "周一",
	time.Tuesday:   "周二",
	time.Wednesday: "周三",
	time.Thursday:  "周四",
	time.Friday:    "周五",
	time.Saturday:  "周六",
}
