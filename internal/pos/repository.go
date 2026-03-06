package pos

import "gorm.io/gorm"

// CreateTransaction inserts a transaction. Accepts db or tx. Returns error.
func CreateTransaction(db *gorm.DB, t *Transaction) error {
	return db.Create(t).Error
}

// CreateTransactionItem inserts a transaction item. Accepts db or tx.
func CreateTransactionItem(db *gorm.DB, item *TransactionItem) error {
	return db.Create(item).Error
}

// CreatePayment inserts a payment. Accepts db or tx.
func CreatePayment(db *gorm.DB, p *Payment) error {
	return db.Create(p).Error
}
