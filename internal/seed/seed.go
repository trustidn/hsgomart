package seed

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type productDef struct {
	Name      string
	SKU       string
	Barcode   string
	CostPrice float64
	SellPrice float64
	Category  int // index into categories slice
	Unit      string
	InitStock int
}

func SeedTenantData(db *gorm.DB, tenantID, ownerID string) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// --- Cashier user ---
	shortID := tenantID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}
	cashierEmail := fmt.Sprintf("kasir-%s@demo.test", shortID)
	hash, err := bcrypt.GenerateFromPassword([]byte("Kasir123!"), 10)
	if err != nil {
		tx.Rollback()
		return err
	}
	var cashierID string
	if err := tx.Raw(`INSERT INTO users (tenant_id, name, email, password_hash, role, status)
		VALUES (?, 'Kasir Demo', ?, ?, 'cashier', 'active') RETURNING id`,
		tenantID, cashierEmail, string(hash)).Scan(&cashierID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// --- Categories ---
	catNames := []string{"Makanan", "Minuman", "Kebutuhan Rumah"}
	catIDs := make([]string, len(catNames))
	for i, name := range catNames {
		if err := tx.Raw(`INSERT INTO categories (tenant_id, name) VALUES (?, ?) RETURNING id`,
			tenantID, name).Scan(&catIDs[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// --- Products ---
	products := []productDef{
		{"Indomie Goreng", "MKN001", "8886008101053", 2500, 3500, 0, "pcs", 48},
		{"Indomie Soto", "MKN002", "8886008101060", 2500, 3500, 0, "pcs", 36},
		{"Chitato Original 68g", "MKN003", "8886013200109", 8000, 11000, 0, "pcs", 24},
		{"Roma Kelapa", "MKN004", "8996001301011", 3000, 5000, 0, "pcs", 20},
		{"Aqua 600ml", "MNM001", "8998009010019", 2000, 3500, 1, "pcs", 60},
		{"Teh Botol Sosro 450ml", "MNM002", "8886008600068", 3500, 5000, 1, "pcs", 36},
		{"Coca Cola 390ml", "MNM003", "8888166390119", 5000, 7000, 1, "pcs", 24},
		{"Good Day Cappuccino 250ml", "MNM004", "8807057088832", 4000, 6000, 1, "pcs", 30},
		{"Rinso Cair 800ml", "RMH001", "8886011213213", 15000, 22000, 2, "pcs", 12},
		{"Molto Pewangi 800ml", "RMH002", "8886011236038", 12000, 18000, 2, "pcs", 10},
		{"Tissue Paseo 250 sheets", "RMH003", "8993053600013", 8000, 12000, 2, "pcs", 15},
		{"Sabun Lifebuoy 100g", "RMH004", "8886008100018", 3500, 5500, 2, "pcs", 20},
	}

	barcodePrefix := shortID[:4]
	productIDs := make([]string, len(products))
	for i, p := range products {
		if err := tx.Raw(`INSERT INTO products (tenant_id, category_id, name, sku, cost_price, sell_price, unit, low_stock_threshold, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, 10, 'active') RETURNING id`,
			tenantID, catIDs[p.Category], p.Name, p.SKU, p.CostPrice, p.SellPrice, p.Unit).Scan(&productIDs[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
		uniqueBarcode := fmt.Sprintf("%s%s", barcodePrefix, p.Barcode)
		if err := tx.Exec(`INSERT INTO product_barcodes (product_id, barcode) VALUES (?, ?)`,
			productIDs[i], uniqueBarcode).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// --- Purchases & inventory ---
	now := time.Now()
	purchaseDate1 := now.Add(-72 * time.Hour)
	purchaseDate2 := now.Add(-48 * time.Hour)

	type purchaseBatch struct {
		productIdx int
		qty        int
		costPrice  float64
	}

	purchaseSets := []struct {
		supplier string
		invoice  string
		date     time.Time
		items    []purchaseBatch
	}{
		{
			supplier: "PT Indofood Sukses Makmur",
			invoice:  "INV-DEMO-001",
			date:     purchaseDate1,
			items: []purchaseBatch{
				{0, 48, 2500}, {1, 36, 2500}, {2, 24, 8000}, {3, 20, 3000},
				{4, 60, 2000}, {5, 36, 3500}, {7, 30, 4000},
			},
		},
		{
			supplier: "PT Unilever Indonesia",
			invoice:  "INV-DEMO-002",
			date:     purchaseDate2,
			items: []purchaseBatch{
				{6, 24, 5000}, {8, 12, 15000}, {9, 10, 12000},
				{10, 15, 8000}, {11, 20, 3500},
			},
		},
	}

	for _, ps := range purchaseSets {
		totalAmount := 0.0
		for _, item := range ps.items {
			totalAmount += float64(item.qty) * item.costPrice
		}

		var purchaseID string
		if err := tx.Raw(`INSERT INTO purchases (tenant_id, supplier_name, invoice_number, total_amount, created_at)
			VALUES (?, ?, ?, ?, ?) RETURNING id`,
			tenantID, ps.supplier, ps.invoice, totalAmount, ps.date).Scan(&purchaseID).Error; err != nil {
			tx.Rollback()
			return err
		}

		for _, item := range ps.items {
			prodID := productIDs[item.productIdx]
			subtotal := float64(item.qty) * item.costPrice

			var purchaseItemID string
			if err := tx.Raw(`INSERT INTO purchase_items (purchase_id, product_id, quantity, cost_price, subtotal)
				VALUES (?, ?, ?, ?, ?) RETURNING id`,
				purchaseID, prodID, item.qty, item.costPrice, subtotal).Scan(&purchaseItemID).Error; err != nil {
				tx.Rollback()
				return err
			}

			if err := tx.Exec(`INSERT INTO inventory_batches (product_id, purchase_item_id, quantity, remaining_quantity, cost_price, created_at)
				VALUES (?, ?, ?, ?, ?, ?)`,
				prodID, purchaseItemID, item.qty, item.qty, item.costPrice, ps.date).Error; err != nil {
				tx.Rollback()
				return err
			}

			if err := tx.Exec(`INSERT INTO inventories (tenant_id, product_id, stock)
				VALUES (?, ?, ?)
				ON CONFLICT (tenant_id, product_id) DO UPDATE SET stock = inventories.stock + EXCLUDED.stock`,
				tenantID, prodID, item.qty).Error; err != nil {
				tx.Rollback()
				return err
			}

			if err := tx.Exec(`INSERT INTO stock_movements (tenant_id, product_id, type, quantity, reference, reference_id, created_at)
				VALUES (?, ?, 'purchase', ?, ?, ?, ?)`,
				tenantID, prodID, item.qty, ps.invoice, purchaseID, ps.date).Error; err != nil {
				tx.Rollback()
				return err
			}

			tx.Exec(`UPDATE products SET last_purchase_price = ? WHERE id = ?`, item.costPrice, prodID)
		}
	}

	// --- Sample transactions (sales) ---
	type saleLine struct {
		productIdx int
		qty        int
	}
	sales := []struct {
		lines  []saleLine
		method string
		offset time.Duration
	}{
		{
			lines:  []saleLine{{0, 3}, {4, 2}, {2, 1}},
			method: "cash",
			offset: -36 * time.Hour,
		},
		{
			lines:  []saleLine{{5, 2}, {7, 1}, {10, 1}},
			method: "cash",
			offset: -24 * time.Hour,
		},
		{
			lines:  []saleLine{{8, 1}, {11, 2}, {6, 2}},
			method: "qris",
			offset: -12 * time.Hour,
		},
		{
			lines:  []saleLine{{1, 5}, {3, 2}, {9, 1}},
			method: "cash",
			offset: -6 * time.Hour,
		},
	}

	for _, sale := range sales {
		saleTime := now.Add(sale.offset)
		totalAmount := 0.0
		for _, line := range sale.lines {
			totalAmount += float64(line.qty) * products[line.productIdx].SellPrice
		}

		var txID string
		if err := tx.Raw(`INSERT INTO transactions (tenant_id, user_id, total_amount, discount_amount, status, created_at)
			VALUES (?, ?, ?, 0, 'completed', ?) RETURNING id`,
			tenantID, ownerID, totalAmount, saleTime).Scan(&txID).Error; err != nil {
			tx.Rollback()
			return err
		}

		for _, line := range sale.lines {
			p := products[line.productIdx]
			prodID := productIDs[line.productIdx]
			subtotal := float64(line.qty) * p.SellPrice
			cogs := float64(line.qty) * p.CostPrice

			tx.Exec(`INSERT INTO transaction_items (transaction_id, product_id, price, quantity, subtotal, unit_cost, cogs)
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				txID, prodID, p.SellPrice, line.qty, subtotal, p.CostPrice, cogs)

			tx.Exec(`UPDATE inventories SET stock = stock - ? WHERE tenant_id = ? AND product_id = ?`,
				line.qty, tenantID, prodID)

			tx.Exec(`UPDATE inventory_batches SET remaining_quantity = remaining_quantity - ?
				WHERE id = (SELECT id FROM inventory_batches WHERE product_id = ? AND remaining_quantity > 0 ORDER BY created_at ASC LIMIT 1)`,
				line.qty, prodID)

			tx.Exec(`INSERT INTO stock_movements (tenant_id, product_id, type, quantity, reference, reference_id, created_at)
				VALUES (?, ?, 'sale', ?, 'POS Sale', ?, ?)`,
				tenantID, prodID, line.qty, txID, saleTime)
		}

		tx.Exec(`INSERT INTO payments (transaction_id, method, amount, created_at)
			VALUES (?, ?, ?, ?)`, txID, sale.method, totalAmount, saleTime)
	}

	return tx.Commit().Error
}
