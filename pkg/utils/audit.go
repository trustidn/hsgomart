package utils

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID         string          `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID   *string         `gorm:"type:uuid"`
	UserID     *string         `gorm:"type:uuid"`
	Action     string          `gorm:"type:varchar(100);not null"`
	EntityType string          `gorm:"type:varchar(100)"`
	EntityID   string          `gorm:"type:varchar(100)"`
	Details    json.RawMessage `gorm:"type:jsonb"`
	CreatedAt  time.Time       `gorm:"autoCreateTime"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

func LogAudit(db *gorm.DB, tenantID, userID, action, entityType, entityID string, details interface{}) {
	var detailsJSON json.RawMessage
	if details != nil {
		if b, err := json.Marshal(details); err == nil {
			detailsJSON = b
		}
	}
	var tid *string
	if tenantID != "" {
		tid = &tenantID
	}
	var uid *string
	if userID != "" {
		uid = &userID
	}
	entry := AuditLog{
		TenantID:   tid,
		UserID:     uid,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    detailsJSON,
	}
	_ = db.Create(&entry).Error
}
