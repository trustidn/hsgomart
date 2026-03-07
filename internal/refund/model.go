package refund

import "time"

type Refund struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID      string    `gorm:"type:uuid;not null;index"`
	TransactionID string    `gorm:"type:uuid;not null;index"`
	UserID        string    `gorm:"type:uuid;not null"`
	Amount        float64   `gorm:"type:numeric;not null"`
	Reason        string    `gorm:"type:text"`
	Status        string    `gorm:"type:varchar(50);default:completed"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (Refund) TableName() string {
	return "refunds"
}
