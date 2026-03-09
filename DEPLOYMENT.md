# Panduan Deployment Production - HSGomart

Panduan lengkap untuk mendeploy HSGomart ke production dari server kosong hingga aplikasi berjalan di **https://hsgomart.com**.

---

## Ringkasan Arsitektur

- **Backend**: Go (Gin), binary `hsmart-server`
- **Frontend**: Vue 3 + Vite (SPA), di-build dan disajikan oleh Go
- **Database**: PostgreSQL 16
- **Migration**: golang-migrate (SQL files di `migrations/`) — **bukan GORM AutoMigrate**
- **Seed**: Perlu dijalankan manual untuk superadmin (`cmd/seed`)

---

## 1. Persiapan Server Ubuntu (22.04/24.04)

```bash
# Update sistem
sudo apt update && sudo apt upgrade -y

# Set timezone
sudo timedatectl set-timezone Asia/Jakarta

# Buat user deploy (opsional, untuk keamanan)
sudo adduser deploy --disabled-password
sudo usermod -aG sudo deploy
```

---

## 2. Install Dependency Server

```bash
sudo apt install -y \
  curl \
  wget \
  git \
  unzip \
  ca-certificates \
  gnupg \
  lsb-release
```

---

## 3. Install Docker dan Docker Compose

```bash
# Tambah Docker GPG key dan repository
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine dan Compose
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Verifikasi
docker --version
docker compose version

# (Opsional) Tambah user ke group docker agar tidak perlu sudo
sudo usermod -aG docker $USER
# Logout lalu login kembali agar group aktif
```

---

## 4. Konfigurasi Firewall

```bash
# Aktifkan UFW
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw --force enable
sudo ufw status
```

---

## 5. Clone Repository

```bash
# Direktori aplikasi (sesuaikan jika perlu)
export APP_DIR=/opt/hsgomart
sudo mkdir -p $APP_DIR
sudo chown $USER:$USER $APP_DIR

# Clone
git clone https://github.com/trustidn/hsgomart.git $APP_DIR
cd $APP_DIR

# Jika deploy branch tertentu
# git checkout main
# git pull origin main
```

---

## 6. Konfigurasi Environment (.env)

```bash
cd $APP_DIR

# Generate JWT_SECRET (WAJIB)
JWT_SECRET=$(openssl rand -hex 32)
echo "JWT_SECRET yang di-generate: $JWT_SECRET"

# Buat file .env
cat > .env << 'ENVEOF'
# Application
APP_ENV=production
APP_PORT=8080

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=hsmart_user
DB_PASSWORD=P4ss#word90901
DB_NAME=hsmart_saas

# JWT - WAJIB (gunakan output dari: openssl rand -hex 32)
JWT_SECRET=PASTE_HASIL_OPENSSL_RAND_HEX_32_DISINI

# CORS - domain production
CORS_ORIGINS=https://hsgomart.com,https://www.hsgomart.com
ENVEOF

# Ganti JWT_SECRET dengan nilai yang sudah di-generate
sed -i "s|JWT_SECRET=.*|JWT_SECRET=$JWT_SECRET|" .env

# Pastikan permission
chmod 600 .env
```

**Catatan:** Jika password mengandung karakter khusus (`#`, `@`, `%`), gunakan tanda kutip di `.env`:
```
DB_PASSWORD="P4ss#word90901"
```

**Generate JWT_SECRET:**
```bash
openssl rand -hex 32
```
Salin output dan paste ke nilai `JWT_SECRET` di `.env`.

---

## 7. Build Docker Image

```bash
cd $APP_DIR

# Build image
docker compose build --no-cache

# Verifikasi image
docker images | grep hsgomart
```

---

## 8. Menjalankan Database

Database PostgreSQL akan start otomatis bersama stack. Untuk menjalankan hanya database terlebih dahulu (opsional, untuk testing):

```bash
cd $APP_DIR
docker compose up -d postgres

# Tunggu hingga healthy
docker compose ps
# Status postgres harus "healthy"
```

---

## 9. Menjalankan Migration

**Mekanisme migration di project ini:**

- Menggunakan **golang-migrate** (SQL files di `migrations/`)
- **Bukan GORM AutoMigrate** — tidak ada `db.AutoMigrate()` di kode
- Migration dijalankan otomatis oleh service `migrate` saat `docker compose up`

**Alur saat `docker compose up`:**
1. Service `postgres` start → tunggu healthcheck OK
2. Service `migrate` run → jalankan `migrate -path /migrations -database ... up`
3. Setelah migrate selesai, service `app` start

**Manual migration (jika perlu run ulang atau troubleshooting):**

```bash
cd $APP_DIR

# Set DATABASE_URL (untuk script migrate.sh)
export DATABASE_URL="postgres://hsmart_user:P4ss#word90901@localhost:5432/hsmart_saas?sslmode=disable"

# Via migrate tool (install dulu: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
./scripts/migrate.sh up

# Atau via Docker
docker compose run --rm migrate
```

---

## 10. Menjalankan Seed Data

**Seed superadmin** harus dijalankan manual setelah migration. Binary `hsmart-seed` sudah termasuk di image.

```bash
cd $APP_DIR

# Pastikan postgres dan migrate sudah selesai
docker compose up -d postgres
sleep 15
docker compose run --rm migrate  # jalankan migration jika belum

# Jalankan seed (binary sudah ada di image)
docker compose run --rm --no-deps \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=hsmart_user \
  -e DB_PASSWORD='P4ss#word90901' \
  -e DB_NAME=hsmart_saas \
  app ./hsmart-seed

# Output: "Superadmin seeded successfully."
```

**Default superadmin:**
- Email: `admin@hsmart.io`
- Password: `SuperAdmin123!`
- Role: `superadmin`

Ubah via env: `SUPERADMIN_EMAIL`, `SUPERADMIN_PASSWORD`.

---

## 11. Menjalankan Aplikasi

```bash
cd $APP_DIR

# Start seluruh stack (postgres + migrate + app)
docker compose up -d

# Verifikasi
docker compose ps

# Log aplikasi
docker compose logs -f app
```

---

## 12. Verifikasi Health Endpoint

```bash
# Dari server
curl -s http://localhost:8080/health | jq

# Expected:
# { "status": "ok", "database": "connected" }
```

---

## 13. Setup Nginx Reverse Proxy

```bash
# Install Nginx
sudo apt install -y nginx

# Buat konfigurasi
sudo tee /etc/nginx/sites-available/hsgomart.com << 'NGINXEOF'
server {
    listen 80;
    server_name hsgomart.com www.hsgomart.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_read_timeout 300s;
        proxy_connect_timeout 75s;
    }

    client_max_body_size 10M;
}
NGINXEOF

# Aktifkan site
sudo ln -sf /etc/nginx/sites-available/hsgomart.com /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default 2>/dev/null

# Test konfigurasi
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

---

## 14. Setup SSL Let's Encrypt

```bash
# Install Certbot
sudo apt install -y certbot python3-certbot-nginx

# Pastikan DNS hsgomart.com dan www.hsgomart.com sudah mengarah ke IP server

# Generate sertifikat
sudo certbot --nginx -d hsgomart.com -d www.hsgomart.com

# Ikuti prompt (email, agree ToS)
# Pilih redirect HTTP ke HTTPS (opsi 2)

# Verifikasi auto-renewal
sudo certbot renew --dry-run
```

---

## 15. Setup Backup Database

```bash
# Buat direktori scripts jika belum ada
mkdir -p $APP_DIR/scripts

# Buat script backup
tee $APP_DIR/scripts/backup-db.sh << 'BACKUPEOF'
#!/bin/bash
set -e
BACKUP_DIR="/opt/hsgomart/backups"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p "$BACKUP_DIR"

docker exec hsmart-postgres pg_dump -U hsmart_user hsmart_saas | gzip > "$BACKUP_DIR/hsmart_${DATE}.sql.gz"

# Hapus backup older than 7 days
find "$BACKUP_DIR" -name "hsmart_*.sql.gz" -mtime +7 -delete
BACKUPEOF

chmod +x $APP_DIR/scripts/backup-db.sh

# Jadwalkan cron (setiap hari jam 02:00)
(crontab -l 2>/dev/null; echo "0 2 * * * $APP_DIR/scripts/backup-db.sh") | crontab -

# Restore dari backup (jika perlu)
# gunzip -c backups/hsmart_YYYYMMDD_HHMMSS.sql.gz | docker exec -i hsmart-postgres psql -U hsmart_user hsmart_saas
```

---

## 16. Cara Update Aplikasi di Production

```bash
cd /opt/hsgomart

# Pull perubahan terbaru
git pull origin main

# Rebuild image dan restart
docker compose build --no-cache app
docker compose up -d

# Migration otomatis berjalan saat stack naik (service migrate)

# Verifikasi
docker compose ps
curl -s http://localhost:8080/health | jq
```

---

## Checklist Production

- [ ] `.env` dikonfigurasi dengan `JWT_SECRET` minimal 32 karakter
- [ ] `DB_PASSWORD` kuat dan aman
- [ ] `CORS_ORIGINS` berisi domain production (https://hsgomart.com, https://www.hsgomart.com)
- [ ] Firewall (UFW) aktif, port 80/443 terbuka
- [ ] SSL aktif via Certbot
- [ ] Backup database dijadwalkan (cron)
- [ ] Seed superadmin sudah dijalankan
- [ ] Health endpoint merespons OK

---

## Troubleshooting

**App tidak start (JWT_SECRET required):**
```
Error: JWT_SECRET is required
```
→ Pastikan `JWT_SECRET` ada di `.env` dan minimal 32 karakter.

**Migration gagal (Dirty database):**
```bash
# Via Docker (ganti <version> dengan angka dari pesan error, misalnya 36)
# Jika password mengandung #, encode sebagai %23
docker compose run --rm migrate -path /migrations \
  -database "postgres://hsmart_user:P4ss%23word90901@postgres:5432/hsmart_saas?sslmode=disable" \
  force <version>

# Atau dari host dengan migrate CLI:
# export DATABASE_URL="postgres://hsmart_user:P4ss%23word90901@localhost:5432/hsmart_saas?sslmode=disable"
# ./scripts/migrate.sh force <version>
```

**Database connection refused:**
- Pastikan service `postgres` healthy: `docker compose ps`
- Cek log: `docker compose logs postgres`

**Port 5432 tidak bisa diakses dari host:**
- docker-compose menggunakan `127.0.0.1:5432` — hanya localhost. Untuk backup, gunakan `docker exec`.
