package tenant

import "time"

type Tenant struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255)"`
	Phone     string    `gorm:"type:varchar(50)"`
	Status    string    `gorm:"type:varchar(50);default:active"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Tenant) TableName() string {
	return "tenants"
}
