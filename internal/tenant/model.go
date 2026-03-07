package tenant

import "time"

type Tenant struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Email       string    `gorm:"type:varchar(255)" json:"email"`
	Phone       string    `gorm:"type:varchar(50)" json:"phone"`
	LogoURL     string    `gorm:"type:text" json:"logo_url"`
	Address     string    `gorm:"type:text" json:"address"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(50);default:active" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Tenant) TableName() string {
	return "tenants"
}
