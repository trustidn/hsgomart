package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID     string         `gorm:"type:uuid;not null;index"`
	Name         string         `gorm:"type:varchar(255)"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex"`
	PasswordHash string         `gorm:"column:password_hash;type:text"`
	Role         string         `gorm:"type:varchar(50)"`
	Status       string         `gorm:"type:varchar(50);default:active"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
