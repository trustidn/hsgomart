package subscription

import "time"

type Plan struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(100)"`
	Price       float64   `gorm:"type:numeric"`
	MaxUsers    int       `gorm:"column:max_users"`
	MaxProducts int       `gorm:"column:max_products"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (Plan) TableName() string {
	return "plans"
}

type Subscription struct {
	ID        string     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID  string     `gorm:"type:uuid;not null;index"`
	PlanID    int        `gorm:"not null"`
	Status    string     `gorm:"type:varchar(50)"`
	StartDate *time.Time `gorm:"column:start_date"`
	EndDate   *time.Time `gorm:"column:end_date"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
