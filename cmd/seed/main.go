package main

import (
	"fmt"
	"log"
	"os"

	"github.com/trustidn/hsmart-saas/pkg/config"
	"github.com/trustidn/hsmart-saas/pkg/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := seedSuperAdmin(db); err != nil {
		log.Fatalf("failed to seed superadmin: %v", err)
	}
	fmt.Println("Superadmin seeded successfully.")
	fmt.Println("  Email:    admin@hsmart.io")
	fmt.Println("  Password: SuperAdmin123!")
	fmt.Println("  Role:     superadmin")
}

func seedSuperAdmin(db *gorm.DB) error {
	email := "admin@hsmart.io"
	password := "SuperAdmin123!"

	if envEmail := os.Getenv("SUPERADMIN_EMAIL"); envEmail != "" {
		email = envEmail
	}
	if envPass := os.Getenv("SUPERADMIN_PASSWORD"); envPass != "" {
		password = envPass
	}

	// Ensure role constraint allows 'superadmin'
	db.Exec("ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check")
	db.Exec("ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('owner', 'cashier', 'superadmin'))")

	var count int64
	db.Raw("SELECT COUNT(*) FROM users WHERE role = 'superadmin'").Scan(&count)
	if count > 0 {
		fmt.Println("Superadmin already exists, skipping.")
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return fmt.Errorf("bcrypt: %w", err)
	}

	tx := db.Begin()

	var tenantID string
	err = tx.Raw("SELECT id FROM tenants WHERE name = 'System Admin' LIMIT 1").Scan(&tenantID).Error
	if err != nil || tenantID == "" {
		err = tx.Raw(`
			INSERT INTO tenants (name, email, status)
			VALUES ('System Admin', ?, 'active')
			RETURNING id
		`, email).Scan(&tenantID).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("create tenant: %w", err)
		}
	}

	err = tx.Exec(`
		INSERT INTO users (tenant_id, name, email, password_hash, role, status)
		VALUES (?, 'Super Admin', ?, ?, 'superadmin', 'active')
	`, tenantID, email, string(hash)).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("create user: %w", err)
	}

	return tx.Commit().Error
}
