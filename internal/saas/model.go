package saas

import "time"

type Settings struct {
	ID           int       `gorm:"primaryKey" json:"id"`
	SaasName     string    `gorm:"column:saas_name;type:varchar(255)" json:"saas_name"`
	LogoURL      string    `gorm:"column:logo_url;type:text" json:"logo_url"`
	Tagline      string    `gorm:"type:text" json:"tagline"`
	BankName     string    `gorm:"column:bank_name;type:varchar(100)" json:"bank_name"`
	BankAccount  string    `gorm:"column:bank_account;type:varchar(50)" json:"bank_account"`
	BankHolder   string    `gorm:"column:bank_holder;type:varchar(100)" json:"bank_holder"`
	ContactEmail    string    `gorm:"column:contact_email;type:varchar(255)" json:"contact_email"`
	ContactPhone    string    `gorm:"column:contact_phone;type:varchar(50)" json:"contact_phone"`
	WhatsappNumber  string    `gorm:"column:whatsapp_number;type:varchar(50)" json:"whatsapp_number"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Settings) TableName() string {
	return "saas_settings"
}
