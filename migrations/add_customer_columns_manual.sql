-- Tambahkan kolom customer_name dan customer_phone ke tabel transactions
-- Jalankan manual jika migration 0041 belum berhasil dijalankan
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS customer_name VARCHAR(255) DEFAULT '';
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS customer_phone VARCHAR(50) DEFAULT '';
