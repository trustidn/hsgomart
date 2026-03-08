package subscription

import "time"

type Plan struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100)" json:"name"`
	Price       float64   `gorm:"type:numeric" json:"price"`
	MaxUsers    int       `gorm:"column:max_users" json:"max_users"`
	MaxProducts int       `gorm:"column:max_products" json:"max_products"`
	Description string    `gorm:"type:text" json:"description"`
	IsActive    *bool     `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Plan) TableName() string {
	return "plans"
}

type Subscription struct {
	ID        string     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TenantID  string     `gorm:"type:uuid;not null;index" json:"tenant_id"`
	PlanID    int        `gorm:"not null" json:"plan_id"`
	Status    string     `gorm:"type:varchar(50)" json:"status"`
	StartDate *time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate   *time.Time `gorm:"column:end_date" json:"end_date"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
