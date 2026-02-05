package repository

import (
	"time"
	"user-center/internal/model"

	"gorm.io/gorm"
)

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository() *AttendanceRepository {
	return &AttendanceRepository{db: GetDB()}
}

// CreateOrUpdate 创建或更新在班记录
func (r *AttendanceRepository) CreateOrUpdate(record *model.CCAttendance) error {
	// 查找是否存在
	var existing model.CCAttendance
	err := r.db.Where("cc_id = ? AND attendance_date = ?", record.CCID, record.AttendanceDate).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(record).Error
	}
	if err != nil {
		return err
	}
	// 更新
	return r.db.Model(&existing).Updates(map[string]interface{}{
		"status":      record.Status,
		"operator_id": record.OperatorID,
	}).Error
}

// BatchCreateOrUpdate 批量创建或更新
func (r *AttendanceRepository) BatchCreateOrUpdate(records []*model.CCAttendance) error {
	for _, record := range records {
		if err := r.CreateOrUpdate(record); err != nil {
			return err
		}
	}
	return nil
}

// GetByDate 获取指定日期的在班记录
func (r *AttendanceRepository) GetByDate(date time.Time) ([]*model.CCAttendance, error) {
	var records []*model.CCAttendance
	err := r.db.Where("attendance_date = ?", date.Format("2006-01-02")).Find(&records).Error
	return records, err
}

// GetByDateRange 获取日期范围内的在班记录
func (r *AttendanceRepository) GetByDateRange(startDate, endDate time.Time) ([]*model.CCAttendance, error) {
	var records []*model.CCAttendance
	err := r.db.Where("attendance_date >= ? AND attendance_date <= ?",
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Find(&records).Error
	return records, err
}

// GetByCCID 获取CC的在班记录
func (r *AttendanceRepository) GetByCCID(ccID int64, startDate, endDate time.Time) ([]*model.CCAttendance, error) {
	var records []*model.CCAttendance
	err := r.db.Where("cc_id = ? AND attendance_date >= ? AND attendance_date <= ?",
		ccID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("attendance_date ASC").
		Find(&records).Error
	return records, err
}

// GetByCCIDAndDate 获取CC指定日期的在班记录
func (r *AttendanceRepository) GetByCCIDAndDate(ccID int64, date time.Time) (*model.CCAttendance, error) {
	var record model.CCAttendance
	err := r.db.Where("cc_id = ? AND attendance_date = ?", ccID, date.Format("2006-01-02")).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &record, err
}

// CountByStatus 统计指定日期各状态人数
func (r *AttendanceRepository) CountByStatus(date time.Time) (map[string]int64, error) {
	type Result struct {
		Status string
		Count  int64
	}
	var results []Result
	err := r.db.Table("cc_attendance").
		Select("status, COUNT(*) as count").
		Where("attendance_date = ?", date.Format("2006-01-02")).
		Group("status").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}
