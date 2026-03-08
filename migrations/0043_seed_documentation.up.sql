INSERT INTO documentation (title, content, sort_order, is_published) VALUES
('Memulai Menggunakan Sistem', 'Selamat datang di sistem POS kami! Berikut langkah-langkah awal untuk memulai:

1. Setelah mendaftar, Anda akan masuk ke Dashboard utama
2. Atur profil toko Anda di menu Settings (nama toko, logo, alamat)
3. Tambahkan kategori produk terlebih dahulu
4. Tambahkan produk beserta harga jual dan harga beli
5. Setelah produk tersedia, Anda bisa mulai bertransaksi melalui POS

Tips: Gunakan fitur trial untuk mencoba semua fitur tanpa biaya!', 1, true),

('Kategori & Produk', 'MEMBUAT KATEGORI
Buka menu Categories untuk membuat kategori produk. Contoh: Makanan, Minuman, Snack, Kebutuhan Rumah Tangga.

MENAMBAH PRODUK
1. Buka menu Products > klik "Add Product"
2. Isi nama produk, pilih kategori, masukkan harga beli & harga jual
3. (Opsional) Tambahkan barcode jika produk memiliki barcode fisik
4. Anda bisa menambahkan beberapa barcode untuk satu produk

BARCODE
- Barcode bisa di-scan langsung di POS untuk mempercepat transaksi
- Satu produk bisa memiliki beberapa barcode (barcode ganda)
- Tambah barcode melalui halaman detail produk', 2, true),

('Menggunakan POS (Point of Sale)', 'CARA MELAKUKAN TRANSAKSI

1. Buka menu POS dari sidebar atau dashboard
2. Cari produk dengan mengetik nama atau scan barcode
3. Klik produk untuk menambahkan ke keranjang belanja
4. Atur jumlah item sesuai kebutuhan
5. Pilih metode pembayaran (Tunai, QRIS, Transfer, E-Wallet, dll)
6. Masukkan jumlah bayar
7. (Opsional) Masukkan nama dan nomor HP pelanggan untuk pengiriman struk via WhatsApp
8. Klik "Bayar" untuk menyelesaikan transaksi

MENGIRIM STRUK VIA WHATSAPP
- Jika nomor HP pelanggan diisi saat checkout, struk akan dikirim langsung ke WhatsApp pelanggan
- Anda juga bisa mencetak struk atau membagikannya secara manual

METODE PEMBAYARAN
Sistem mendukung berbagai metode pembayaran: Tunai, QRIS, Transfer Bank, E-Wallet (GoPay, OVO, Dana, dll), dan Kartu Debit/Kredit.', 3, true),

('Inventory & Stok', 'MELIHAT STOK
Buka menu Inventory untuk melihat stok seluruh produk. Produk dengan stok rendah akan diberi peringatan.

MENAMBAH STOK MELALUI PURCHASES
1. Buka menu Purchases > klik "New Purchase"
2. Pilih produk yang ingin ditambah stoknya
3. Masukkan jumlah, harga beli per unit, dan (opsional) tanggal kadaluarsa
4. Simpan - stok akan otomatis bertambah

RIWAYAT PERGERAKAN STOK
Menu Inventory History mencatat semua pergerakan stok:
- Penambahan dari Purchase
- Pengurangan dari penjualan (POS)
- Penyesuaian manual
- Koreksi dari Stock Opname

STOCK OPNAME
Stock Opname digunakan untuk mencocokkan stok fisik dengan stok di sistem:
1. Buat Stock Opname baru
2. Hitung stok fisik setiap produk
3. Input hasil perhitungan
4. Sistem akan menampilkan selisih antara stok sistem dan stok fisik
5. Approve untuk menyesuaikan stok sistem secara otomatis', 4, true),

('Shift Kasir', 'FUNGSI SHIFT
Shift digunakan untuk memisahkan waktu kerja kasir dan menghitung total transaksi per shift.

CARA MENGGUNAKAN
1. Kasir membuka shift di awal jam kerja melalui menu Shifts > Open Shift
2. Selama shift berlangsung, semua transaksi tercatat dalam shift tersebut
3. Di akhir jam kerja, kasir menutup shift (Close Shift)
4. Sistem menampilkan ringkasan: total transaksi, total penjualan, dan detail pembayaran

MANFAAT
- Memudahkan perhitungan uang di kasir saat pergantian shift
- Owner dapat melihat performa kasir per shift
- Laporan shift tersedia di menu Reports', 5, true),

('Laporan & Analisis', 'Menu Reports menyediakan berbagai laporan untuk menganalisis bisnis Anda:

LAPORAN PENJUALAN
- Penjualan harian dan tren
- Detail per transaksi (termasuk nama pelanggan & no HP)
- Perbandingan dengan periode sebelumnya

LAPORAN PRODUK
- Produk terlaris
- Margin keuntungan per produk
- Produk yang jarang terjual

LAPORAN KASIR
- Performa penjualan per kasir
- Rekap shift kasir

LAPORAN PEMBAYARAN
- Breakdown metode pembayaran
- Total per metode (tunai, QRIS, transfer, dll)

EXPORT DATA
Laporan dapat di-export ke Excel (XLSX) atau PDF untuk keperluan pembukuan.', 6, true);
