package refund

import "gorm.io/gorm"

func CreateRefund(db *gorm.DB, r *Refund) error {
	return db.Create(r).Error
}

func ListRefundsByTenant(db *gorm.DB, tenantID string, limit, offset int) ([]Refund, error) {
	var list []Refund
	q := db.Where("tenant_id = ?", tenantID).Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	err := q.Find(&list).Error
	return list, err
}
