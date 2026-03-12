# Alur Kerja Cetak Struk – HSGoMart

Dokumen ini menjelaskan alur kerja lengkap fitur cetak struk pada sistem HSGoMart, yang berjalan baik di desktop maupun mobile. Dapat digunakan sebagai panduan implementasi di sistem lain.

---

## 1. Ringkasan Arsitektur

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  POS Checkout   │────▶│  Receipt Modal   │────▶│ Output: Print   │
│  (cart → API)   │     │  (tampil struk)  │     │ PDF / WhatsApp  │
└─────────────────┘     └──────────────────┘     └─────────────────┘
```

**Prinsip utama:**
- Struk dirender sebagai **HTML** dengan komponen reusable
- Cetak memakai **`window.print()`** browser (print-to-printer atau save as PDF)
- PDF ekspor memakai **jsPDF** di client
- WhatsApp memakai **wa.me** redirect (text pre-filled)

---

## 2. Alur Kerja (Flow Diagram)

```
[User selesai belanja] 
        │
        ▼
[Klik Checkout] ──▶ [Isi metode pembayaran, nominal, (opsional) nama & HP pelanggan]
        │
        ▼
[POST /api/pos/checkout] ──▶ Backend: validasi, kurangi stok, buat transaction & payments
        │
        ▼
[Response: transaction_id, total, change, cashier]
        │
        ▼
[Receipt Modal muncul] ──▶ Tampilkan komponen <Receipt> dengan data dari response + cart
        │
        ├──▶ [Print]     ──▶ window.print() ──▶ Browser native print dialog
        │
        ├──▶ [PDF]      ──▶ jsPDF.generate() ──▶ Download file .pdf
        │
        └──▶ [WhatsApp] ──▶ wa.me/?text=... atau wa.me/{phone}?text=...
```

---

## 3. Komponen Kunci

### 3.1 Komponen Struk (`Receipt.vue`)

**Fungsi:** Menampilkan struk dalam format yang siap dicetak/dishare.

**Props:**
- `storeName` – Nama toko
- `date` – Tanggal transaksi
- `transactionId` – ID transaksi
- `cashier` – Nama kasir
- `items` – Array: `{ name, quantity, price }`
- `total`, `paidAmount`, `change`

**Elemen penting:**
- Container dengan `id="receipt-print"` untuk CSS print
- Class `print:bg-white print:text-black` agar di print selalu hitam-putih
- Font monospace untuk tampilan mirip struk thermal

### 3.2 CSS untuk Print

```css
@media print {
  body * { visibility: hidden; }
  #receipt-print, #receipt-print * { visibility: visible; }
  #receipt-print {
    position: absolute;
    left: 0; top: 0;
    width: 100%;
    max-width: none;
    box-shadow: none;
  }
}
```

**Poin penting:**
1. **Sembunyikan semua elemen** kecuali `#receipt-print` dan anaknya
2. **Tidak ada overlay/modal** yang tercetak
3. **Lebar penuh** supaya pas untuk thermal printer 80mm atau A4

### 3.3 Utilitas PDF (`receipt-pdf.js`)

- **`generateReceiptPDF(data)`** – Buat PDF 80x200mm (format struk thermal) dengan jsPDF
- **`buildReceiptText(data)`** – Format struk sebagai teks untuk WhatsApp (Markdown bold untuk nama toko)

---

## 4. Integrasi dengan POS

### 4.1 Setelah Checkout Berhasil

```javascript
// Data struk = gabungan response API + data cart
receiptData.value = {
  storeName: tenantStore.storeName(),
  date: new Date(),
  transactionId: result?.transaction_id ?? '',
  cashier: result?.cashier ?? '',
  items: cartItems.value.map(i => ({ ...i })),
  total: result?.total ?? totalAmount.value,
  paidAmount: checkoutForm.value.paid_amount,
  change: result?.change ?? 0,
  customerName: checkoutForm.value.customer_name?.trim() || '',
  customerPhone: checkoutForm.value.customer_phone?.trim() || '',
}
showReceiptModal.value = true
```

### 4.2 Tombol Aksi di Modal

| Tombol    | Aksi                        | Implementasi                          |
|----------|-----------------------------|----------------------------------------|
| **Print** | Cetak langsung / Save PDF   | `window.print()`                       |
| **PDF**   | Unduh file PDF              | `generateReceiptPDF(data).save(...)`  |
| **WhatsApp** | Kirim via WA             | `wa.me/{phone}?text={encoded}`        |
| **Selesai** | Tutup modal, kosongkan cart | Reset state                           |

### 4.3 Logika WhatsApp

```javascript
function shareWhatsApp() {
  const text = encodeURIComponent(buildReceiptText(receiptData.value))
  const phone = (receiptData.value.customerPhone || '').trim().replace(/\D/g, '')
  const waPhone = phone ? (phone.startsWith('62') ? phone : '62' + phone.replace(/^0+/, '')) : ''
  const url = waPhone
    ? `https://wa.me/${waPhone}?text=${text}`   // Langsung ke chat pelanggan
    : `https://wa.me/?text=${text}`              // Buka WA, user pilih contact
  window.open(url, '_blank')
}
```

- Jika `customerPhone` diisi di checkout: buka chat ke nomor tersebut
- Jika kosong: buka WhatsApp dengan teks struk siap dikirim

---

## 5. Mobile vs Desktop

| Aspek            | Desktop                         | Mobile                               |
|------------------|----------------------------------|--------------------------------------|
| Print            | `window.print()`                 | Same (OS print dialog / share)       |
| PDF              | Download                         | Save to files / share               |
| WhatsApp         | `window.open()` di tab baru      | Buka WhatsApp app langsung          |
| Layout modal     | Lebar tetap, max-w-sm            | Responsif, max-h-[90vh]             |
| Layout struk     | Sama, max-w-sm                   | Sama                                 |

**Tidak ada perbedaan kode** untuk cetak antara mobile dan desktop; browser mengurus dialog print dan integrasi dengan printer/aplikasi.

---

## 6. API Backend

### 6.1 Checkout (buat transaksi)

```
POST /api/pos/checkout
Body: {
  items: [{ product_id, quantity }],
  payment_method: "cash" | "card" | "qris" | "ewallet" | "transfer",
  paid_amount: number,
  customer_name?: string,
  customer_phone?: string  // untuk pre-fill WhatsApp
}

Response: {
  transaction_id, total, change, cashier
}
```

### 6.2 Get Receipt (untuk view ulang di Reports)

```
GET /api/pos/receipt/:id

Response: {
  transaction: { id, total_amount, created_at, cashier },
  items: [{ product_name, quantity, price, subtotal }],
  payments: [{ method, amount }]
}
```

---

## 7. Checklist Implementasi untuk Sistem Lain

### Frontend

- [ ] Komponen struk reusable dengan `id="receipt-print"`
- [ ] CSS `@media print` yang menyembunyikan semua kecuali `#receipt-print`
- [ ] `window.print()` tanpa custom print server
- [ ] jsPDF untuk export PDF (opsional, format 80x200mm untuk thermal)
- [ ] `buildReceiptText()` untuk teks WhatsApp
- [ ] Redirect `wa.me/?text=...` atau `wa.me/{phone}?text=...`

### Backend

- [ ] Endpoint checkout yang mengembalikan `transaction_id`, `total`, `change`, `cashier`
- [ ] Endpoint GET receipt untuk view ulang transaksi
- [ ] Kolom `customer_name`, `customer_phone` di tabel transaksi (opsional)

### UX

- [ ] Input nomor HP pelanggan di checkout (opsional)
- [ ] Modal struk pasca-checkout dengan tombol Print, PDF, WhatsApp
- [ ] Struk readable di mobile (font, lebar, spacing)

---

## 8. File Referensi

| File                    | Fungsi                                             |
|-------------------------|----------------------------------------------------|
| `frontend/src/components/Receipt.vue` | Komponen struk HTML + CSS print                    |
| `frontend/src/utils/receipt-pdf.js`    | Generate PDF & teks WhatsApp                      |
| `frontend/src/pages/POS.vue`          | Checkout, modal struk, tombol Print/PDF/WhatsApp |
| `frontend/src/pages/Reports.vue`      | View & cetak struk transaksi lama                 |
| `internal/report/handler.go`          | `GetReceipt` API                                  |
| `internal/pos/`                       | Checkout service & handler                         |

---

## 9. Catatan untuk Thermal Printer

Untuk printer thermal 80mm:
- Format PDF: `new jsPDF({ format: [80, 200] })` (mm)
- Font kecil (7–8pt)
- Hindari grafik kompleks
- `window.print()` → pilih printer thermal di dialog browser

---

*Dokumen ini dibuat dari implementasi HSGoMart untuk keperluan reuse di sistem lain.*
