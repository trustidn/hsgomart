package order

import "time"

type SubscriptionOrder struct {
	ID              string     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TenantID        string     `gorm:"type:uuid;not null;index" json:"tenant_id"`
	PlanID          int        `gorm:"not null" json:"plan_id"`
	Amount          float64    `gorm:"type:numeric;not null" json:"amount"`
	Status          string     `gorm:"type:varchar(50);not null;default:pending_payment" json:"status"`
	PaymentProofURL string     `gorm:"type:text" json:"payment_proof_url"`
	InvoiceNumber   string     `gorm:"type:varchar(50);not null" json:"invoice_number"`
	Notes           string     `gorm:"type:text" json:"notes"`
	AdminNotes      string     `gorm:"type:text" json:"admin_notes"`
	ReviewedBy      *string    `gorm:"type:uuid" json:"reviewed_by"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	PaidAt          *time.Time `json:"paid_at"`
	ReviewedAt      *time.Time `json:"reviewed_at"`
}

func (SubscriptionOrder) TableName() string {
	return "subscription_orders"
}
