package repository

import (
	"time"
	"user-center/internal/model"

	"gorm.io/gorm"
)

type LeadAllocationRepository struct {
	db *gorm.DB
}

func NewLeadAllocationRepository() *LeadAllocationRepository {
	return &LeadAllocationRepository{db: GetDB()}
}

// CreateOrUpdate 创建或更新例子分配记录
func (r *LeadAllocationRepository) CreateOrUpdate(record *model.CCLeadAllocation) error {
	var existing model.CCLeadAllocation
	err := r.db.Where("cc_id = ? AND allocation_date = ?", record.CCID, record.AllocationDate.Format("2006-01-02")).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(record).Error
	}
	if err != nil {
		return err
	}
	// 更新
	return r.db.Model(&existing).Updates(map[string]interface{}{
		"expected_allocation": record.ExpectedAllocation,
		"actual_allocation":   record.ActualAllocation,
		"expected_supplement": record.ExpectedSupplement,
		"actual_supplement":   record.ActualSupplement,
		"overdraft":           record.Overdraft,
		"processed_overdraft": record.ProcessedOverdraft,
		"pending_overdraft":   record.PendingOverdraft,
		"is_allocated":        record.IsAllocated,
		"allocation_rule":     record.AllocationRule,
		"allocation_reason":   record.AllocationReason,
	}).Error
}

// GetByDate 获取指定日期的分配记录
func (r *LeadAllocationRepository) GetByDate(date time.Time) ([]*model.CCLeadAllocation, error) {
	var records []*model.CCLeadAllocation
	err := r.db.Where("allocation_date = ?", date.Format("2006-01-02")).Find(&records).Error
	return records, err
}

// List 分页查询分配记录
func (r *LeadAllocationRepository) List(query *model.LeadAllocationQuery, date time.Time) (*model.PageResult, error) {
	var total int64
	var records []*model.CCLeadAllocation

	db := r.db.Model(&model.CCLeadAllocation{}).Where("allocation_date = ?", date.Format("2006-01-02"))

	if query.CCID != nil {
		db = db.Where("cc_id = ?", *query.CCID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (query.PageNum - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Order("cc_id ASC").Find(&records).Error; err != nil {
		return nil, err
	}

	return &model.PageResult{
		Rows:  records,
		Total: total,
	}, nil
}

// GetByCCIDAndDate 获取CC指定日期的分配记录
func (r *LeadAllocationRepository) GetByCCIDAndDate(ccID int64, date time.Time) (*model.CCLeadAllocation, error) {
	var record model.CCLeadAllocation
	err := r.db.Where("cc_id = ? AND allocation_date = ?", ccID, date.Format("2006-01-02")).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &record, err
}

// UpdateActualAllocation 更新实际分配数
func (r *LeadAllocationRepository) UpdateActualAllocation(ccID int64, date time.Time, actual int) error {
	return r.db.Model(&model.CCLeadAllocation{}).
		Where("cc_id = ? AND allocation_date = ?", ccID, date.Format("2006-01-02")).
		Update("actual_allocation", actual).Error
}

// UpdateOverdraft 更新透支数据
func (r *LeadAllocationRepository) UpdateOverdraft(ccID int64, date time.Time, overdraft, processed, pending int) error {
	return r.db.Model(&model.CCLeadAllocation{}).
		Where("cc_id = ? AND allocation_date = ?", ccID, date.Format("2006-01-02")).
		Updates(map[string]interface{}{
			"overdraft":           overdraft,
			"processed_overdraft": processed,
			"pending_overdraft":   pending,
		}).Error
}

// SumByDate 统计指定日期的分配汇总
func (r *LeadAllocationRepository) SumByDate(date time.Time) (map[string]int64, error) {
	type Result struct {
		TotalExpected  int64
		TotalActual    int64
		TotalOverdraft int64
	}
	var result Result
	err := r.db.Table("cc_lead_allocation").
		Select("SUM(expected_allocation) as total_expected, SUM(actual_allocation) as total_actual, SUM(overdraft) as total_overdraft").
		Where("allocation_date = ?", date.Format("2006-01-02")).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return map[string]int64{
		"totalExpected":  result.TotalExpected,
		"totalActual":    result.TotalActual,
		"totalOverdraft": result.TotalOverdraft,
	}, nil
}

// GetByCCID 获取CC在日期范围内的分配记录
func (r *LeadAllocationRepository) GetByCCID(ccID int64, startDate, endDate time.Time) ([]*model.CCLeadAllocation, error) {
	var records []*model.CCLeadAllocation
	err := r.db.Where("cc_id = ? AND allocation_date >= ? AND allocation_date <= ?",
		ccID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("allocation_date DESC").Find(&records).Error
	return records, err
}
