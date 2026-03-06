package user

import (
	"gorm.io/gorm"
)

// CreateUser inserts a user. Caller must set TenantID and hashed password.
func CreateUser(db *gorm.DB, u *User) error {
	return db.Create(u).Error
}

// FindUsersByTenant returns all users for the tenant (tenant isolation).
func FindUsersByTenant(db *gorm.DB, tenantID string) ([]User, error) {
	var users []User
	err := db.Where("tenant_id = ?", tenantID).Find(&users).Error
	return users, err
}

// FindUserByID returns one user by ID scoped by tenant_id. Returns gorm.ErrRecordNotFound if not found.
func FindUserByID(db *gorm.DB, tenantID, userID string) (*User, error) {
	var u User
	err := db.Where("tenant_id = ? AND id = ?", tenantID, userID).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateUser updates fields for the user scoped by tenant_id. Returns gorm.ErrRecordNotFound if not found.
func UpdateUser(db *gorm.DB, tenantID, userID string, updates map[string]interface{}) error {
	res := db.Model(&User{}).Where("tenant_id = ? AND id = ?", tenantID, userID).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteUser deletes the user scoped by tenant_id. Returns gorm.ErrRecordNotFound if not found.
func DeleteUser(db *gorm.DB, tenantID, userID string) error {
	res := db.Where("tenant_id = ? AND id = ?", tenantID, userID).Delete(&User{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
