UPDATE documentation SET content = '
=============================================
PANDUAN DEPLOYMENT KE SERVER PRODUCTION
=============================================
Server: Ubuntu 22.04/24.04 LTS
Spesifikasi Minimum: 2 CPU, 2-4 GB RAM, 20 GB SSD
=============================================

=== BAGIAN 1: PERSIAPAN SERVER ===

1. UPDATE SISTEM
   sudo apt update && sudo apt upgrade -y
   sudo apt install -y curl wget git ufw nano htop unzip software-properties-common

2. BUAT USER NON-ROOT (jika belum ada)
   adduser deploy
   usermod -aG sudo deploy
   su - deploy

3. KONFIGURASI FIREWALL (UFW)
   sudo ufw allow OpenSSH
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw enable
   sudo ufw status

4. KONFIGURASI SSH (Opsional tapi direkomendasikan)
   - Nonaktifkan login root via SSH:
     sudo nano /etc/ssh/sshd_config
     → Ubah: PermitRootLogin no
     → Ubah: PasswordAuthentication no (jika menggunakan SSH key)
     sudo systemctl restart sshd


=== BAGIAN 2: INSTALL POSTGRESQL ===

1. INSTALL POSTGRESQL 16
   sudo sh -c ''echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list''
   curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg
   sudo apt update
   sudo apt install -y postgresql-16

2. BUAT DATABASE DAN USER
   sudo -u postgres psql
   CREATE USER hsmart_user WITH PASSWORD ''password_anda_yang_kuat'';
   CREATE DATABASE hsmart_saas OWNER hsmart_user;
   GRANT ALL PRIVILEGES ON DATABASE hsmart_saas TO hsmart_user;
   \q

3. KONFIGURASI POSTGRESQL UNTUK PERFORMA
   sudo nano /etc/postgresql/16/main/postgresql.conf

   Ubah parameter berikut (untuk server 2-4 GB RAM):
   shared_buffers = 512MB
   effective_cache_size = 1536MB
   maintenance_work_mem = 128MB
   work_mem = 4MB
   max_connections = 100
   checkpoint_completion_target = 0.9
   wal_buffers = 16MB
   default_statistics_target = 100
   random_page_cost = 1.1
   effective_io_concurrency = 200

   sudo systemctl restart postgresql


=== BAGIAN 3: INSTALL DOCKER & DOCKER COMPOSE ===

1. INSTALL DOCKER
   curl -fsSL https://get.docker.com | sudo sh
   sudo usermod -aG docker deploy
   newgrp docker

2. VERIFIKASI INSTALASI
   docker --version
   docker compose version


=== BAGIAN 4: MENARIK KODE DARI REPOSITORY ===

1. SETUP SSH KEY UNTUK GITHUB (jika repository private)
   ssh-keygen -t ed25519 -C "deploy@server"
   cat ~/.ssh/id_ed25519.pub
   → Tambahkan public key ini ke GitHub repository → Settings → Deploy keys

2. CLONE REPOSITORY
   cd /home/deploy
   git clone git@github.com:username/hsgomart.git
   cd hsgomart

3. SETUP ENVIRONMENT
   cp .env.example .env
   nano .env

   Isi konfigurasi (WAJIB untuk production):
   APP_ENV=production
   APP_PORT=8080
   DB_HOST=host.docker.internal   # atau IP server jika DB di luar Docker
   DB_PORT=5432
   DB_USER=hsmart_user
   DB_PASSWORD=password_anda_yang_kuat
   DB_NAME=hsmart_saas
   JWT_SECRET=<WAJIB minimal 32 karakter, lihat cara generate di bawah>
   CORS_ORIGINS=https://domain-anda.com,https://www.domain-anda.com

   ⚠ PENTING - JWT_SECRET:
   - WAJIB diisi, sistem tidak akan berjalan tanpa JWT_SECRET
   - Di production, MINIMAL 32 karakter
   - Generate JWT secret yang aman:
     openssl rand -hex 32

   ⚠ PENTING - DB_PASSWORD:
   - JANGAN gunakan password default
   - Generate password yang kuat:
     openssl rand -base64 24


=== BAGIAN 5: BUILD DAN JALANKAN APLIKASI ===

1. BUILD DOCKER IMAGE
   docker compose build --no-cache

2. JALANKAN DATABASE MIGRATION
   docker compose up migrate

3. BUAT SUPERADMIN (pertama kali)
   docker compose run --rm app go run ./cmd/seed
   → Catat email dan password superadmin

4. JALANKAN APLIKASI
   docker compose up -d app

5. CEK STATUS
   docker compose ps
   docker compose logs -f app
   curl http://localhost:8080/health

   Aplikasi berjalan di http://localhost:8080

   Catatan: Container berjalan sebagai non-root user untuk keamanan,
   dan memiliki HEALTHCHECK otomatis setiap 30 detik.


=== BAGIAN 6: SETUP DOMAIN DAN SSL ===

1. INSTALL NGINX
   sudo apt install -y nginx

2. KONFIGURASI NGINX SEBAGAI REVERSE PROXY
   sudo nano /etc/nginx/sites-available/hsgomart

   Isi dengan:
   server {
       listen 80;
       server_name domain-anda.com www.domain-anda.com;

       client_max_body_size 10M;

       location / {
           proxy_pass http://127.0.0.1:8080;
           proxy_http_version 1.1;
           proxy_set_header Upgrade $http_upgrade;
           proxy_set_header Connection ''upgrade'';
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
           proxy_cache_bypass $http_upgrade;
       }

       location /uploads/ {
           alias /home/deploy/hsgomart/uploads/;
           expires 30d;
           add_header Cache-Control "public, immutable";
       }
   }

   sudo ln -s /etc/nginx/sites-available/hsgomart /etc/nginx/sites-enabled/
   sudo rm /etc/nginx/sites-enabled/default
   sudo nginx -t
   sudo systemctl reload nginx

3. INSTALL SSL DENGAN CERTBOT (Let''s Encrypt)
   sudo apt install -y certbot python3-certbot-nginx
   sudo certbot --nginx -d domain-anda.com -d www.domain-anda.com

   Certbot akan otomatis:
   - Membuat sertifikat SSL
   - Mengkonfigurasi Nginx untuk HTTPS
   - Redirect HTTP ke HTTPS

4. AUTO-RENEW SSL
   Certbot sudah mengatur auto-renewal. Verifikasi:
   sudo certbot renew --dry-run

   Catatan: Sistem mendeteksi X-Forwarded-Proto dari Nginx untuk
   mengaktifkan HSTS (Strict-Transport-Security) secara otomatis.


=== BAGIAN 7: BACKUP & RESTORE DATABASE ===

1. BACKUP DATABASE (Manual)
   pg_dump -U hsmart_user -h localhost -d hsmart_saas -F c -f backup_$(date +%Y%m%d_%H%M%S).dump

2. RESTORE DATABASE
   pg_restore -U hsmart_user -h localhost -d hsmart_saas -c backup_file.dump

   Atau jika database belum ada:
   createdb -U hsmart_user hsmart_saas_restore
   pg_restore -U hsmart_user -h localhost -d hsmart_saas_restore backup_file.dump

3. BACKUP OTOMATIS (CRON JOB)
   mkdir -p /home/deploy/backups

   Buat script backup:
   nano /home/deploy/backup_db.sh

   Isi:
   #!/bin/bash
   BACKUP_DIR="/home/deploy/backups"
   TIMESTAMP=$(date +%Y%m%d_%H%M%S)
   FILENAME="hsmart_backup_${TIMESTAMP}.dump"

   pg_dump -U hsmart_user -h localhost -d hsmart_saas -F c -f "${BACKUP_DIR}/${FILENAME}"

   # Hapus backup lebih dari 30 hari
   find ${BACKUP_DIR} -name "hsmart_backup_*.dump" -mtime +30 -delete

   echo "Backup selesai: ${FILENAME}"

   chmod +x /home/deploy/backup_db.sh

   Setup cron (backup setiap hari jam 2 pagi):
   crontab -e
   0 2 * * * /home/deploy/backup_db.sh >> /home/deploy/backups/backup.log 2>&1

4. BACKUP KE REMOTE / CLOUD (Opsional)
   Gunakan rclone untuk sync ke Google Drive, S3, atau storage lain:
   sudo apt install -y rclone
   rclone config
   → Ikuti wizard untuk setup remote storage

   Tambahkan ke script backup:
   rclone copy ${BACKUP_DIR}/${FILENAME} remote:hsgomart-backups/

5. TRANSFER BACKUP DARI DEVELOPMENT KE PRODUCTION
   Dari mesin development:
   pg_dump -U hsmart_user -h localhost -d hsmart_saas -F c -f dev_backup.dump
   scp dev_backup.dump deploy@server-ip:/home/deploy/

   Di server production:
   pg_restore -U hsmart_user -h localhost -d hsmart_saas -c /home/deploy/dev_backup.dump


=== BAGIAN 8: HARDENING KEAMANAN ===

1. FAIL2BAN (Proteksi Brute Force)
   sudo apt install -y fail2ban
   sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
   sudo nano /etc/fail2ban/jail.local

   Pastikan:
   [sshd]
   enabled = true
   port = ssh
   maxretry = 5
   bantime = 3600

   sudo systemctl enable fail2ban
   sudo systemctl start fail2ban

2. AUTOMATIC SECURITY UPDATES
   sudo apt install -y unattended-upgrades
   sudo dpkg-reconfigure -plow unattended-upgrades

3. BATASI AKSES POSTGRESQL
   sudo nano /etc/postgresql/16/main/pg_hba.conf

   Pastikan hanya koneksi lokal yang diizinkan:
   local   all   hsmart_user   md5
   host    all   hsmart_user   127.0.0.1/32   md5
   host    all   hsmart_user   172.16.0.0/12  md5   # Docker network

   sudo systemctl restart postgresql

4. SECURITY HEADERS (Sudah Built-in)
   Sistem sudah mengimplementasikan security headers berikut:
   - X-Content-Type-Options: nosniff
   - X-Frame-Options: DENY
   - X-XSS-Protection: 1; mode=block
   - Referrer-Policy: strict-origin-when-cross-origin
   - Permissions-Policy: camera=(), microphone=(), geolocation=()
   - Strict-Transport-Security (otomatis saat HTTPS/reverse proxy)

   Opsional, tambahkan di Nginx:
   add_header Content-Security-Policy "default-src ''self''; img-src ''self'' data:; style-src ''self'' ''unsafe-inline''; script-src ''self'' ''unsafe-inline''" always;

5. NONAKTIFKAN INFORMASI VERSI
   sudo nano /etc/nginx/nginx.conf
   → Tambahkan di block http: server_tokens off;
   sudo systemctl reload nginx

6. RATE LIMITING (Sudah Built-in)
   Sistem sudah memiliki rate limiting pada:
   - Login: 5 request/menit per IP
   - Register: 3 request/menit per IP
   - Refresh Token: 10 request/menit per IP
   - API umum: 60 request/menit per IP
   - Admin API: 30 request/menit per IP

7. KEAMANAN DOCKER
   - Container berjalan sebagai non-root user (appuser)
   - Binary dikompilasi dengan stripped symbols (-ldflags="-s -w")
   - Image minimal (Alpine) mengurangi attack surface
   - HEALTHCHECK otomatis setiap 30 detik


=== BAGIAN 9: OPTIMALISASI PERFORMA ===

1. NGINX CACHING & GZIP
   sudo nano /etc/nginx/nginx.conf

   Tambahkan di block http:
   gzip on;
   gzip_vary on;
   gzip_proxied any;
   gzip_comp_level 6;
   gzip_types text/plain text/css application/json application/javascript text/xml application/xml text/javascript image/svg+xml;

   # Cache static assets
   Tambahkan di block server:
   location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff2?)$ {
       expires 30d;
       add_header Cache-Control "public, immutable";
   }

   sudo systemctl reload nginx

2. DOCKER RESOURCE LIMITS
   Di docker-compose.yml, tambahkan:
   services:
     app:
       deploy:
         resources:
           limits:
             cpus: ''1.5''
             memory: 1536M
           reservations:
             cpus: ''0.5''
             memory: 512M

3. POSTGRESQL CONNECTION POOLING (Opsional, untuk traffic tinggi)
   Install PgBouncer:
   sudo apt install -y pgbouncer
   sudo nano /etc/pgbouncer/pgbouncer.ini

   [databases]
   hsmart_saas = host=127.0.0.1 port=5432 dbname=hsmart_saas

   [pgbouncer]
   listen_addr = 127.0.0.1
   listen_port = 6432
   auth_type = md5
   pool_mode = transaction
   max_client_conn = 200
   default_pool_size = 20

4. FRONTEND OPTIMIZATION (Sudah Built-in)
   - Code splitting: Vue, XLSX, jsPDF di-chunk terpisah
   - Lazy loading routes: Setiap halaman di-load saat dibutuhkan
   - Binary Go dikompilasi dengan stripped debug info

5. MONITORING (Opsional)
   Install Netdata untuk monitoring real-time:
   bash <(curl -Ss https://my-netdata.io/kickstart.sh)
   → Akses di http://server-ip:19999


=== BAGIAN 10: PROSEDUR UPDATE APLIKASI ===

1. TARIK PERUBAHAN TERBARU
   cd /home/deploy/hsgomart
   git pull origin main

2. REBUILD DAN RESTART
   docker compose build --no-cache
   docker compose up -d migrate
   docker compose up -d app

3. VERIFIKASI
   docker compose ps
   docker compose logs --tail=50 app
   curl -s http://localhost:8080/health


=== BAGIAN 11: ENVIRONMENT VARIABLES ===

Daftar lengkap environment variables:

WAJIB:
- APP_ENV        : production (WAJIB untuk mode production)
- JWT_SECRET     : Minimal 32 karakter (WAJIB, gunakan: openssl rand -hex 32)
- DB_PASSWORD    : Password database yang kuat
- DB_USER        : Username database
- DB_NAME        : Nama database
- DB_HOST        : Host database (host.docker.internal untuk Docker)
- DB_PORT        : Port database (default: 5432)

OPSIONAL:
- APP_PORT       : Port aplikasi (default: 8080)
- CORS_ORIGINS   : Comma-separated domain (contoh: https://pos.domain.com)


=== BAGIAN 12: TROUBLESHOOTING ===

MASALAH: Aplikasi tidak bisa start
- JWT_SECRET kosong → Set JWT_SECRET di .env
- JWT_SECRET terlalu pendek → Minimal 32 karakter di production
- Cek log: docker compose logs app

MASALAH: Aplikasi tidak bisa diakses
- Cek status Docker: docker compose ps
- Cek health: curl http://localhost:8080/health
- Cek Nginx: sudo nginx -t && sudo systemctl status nginx
- Cek firewall: sudo ufw status

MASALAH: Database connection error
- Cek PostgreSQL: sudo systemctl status postgresql
- Cek koneksi: psql -U hsmart_user -h localhost -d hsmart_saas
- Cek log: sudo tail -f /var/log/postgresql/postgresql-16-main.log

MASALAH: SSL certificate expired
- Renew manual: sudo certbot renew
- Cek expiry: sudo certbot certificates

MASALAH: Disk penuh
- Cek usage: df -h
- Bersihkan Docker: docker system prune -a
- Cek backup: du -sh /home/deploy/backups/
- Bersihkan log: sudo journalctl --vacuum-time=7d

MASALAH: Rate limiting terlalu ketat
- Rate limit bersifat per-IP dan in-memory
- Restart container akan mereset semua rate limit counter


=== CHECKLIST DEPLOYMENT ===

[ ] Server Ubuntu diupdate
[ ] Firewall (UFW) dikonfigurasi
[ ] PostgreSQL diinstall dan dikonfigurasi
[ ] Docker & Docker Compose diinstall
[ ] Repository di-clone
[ ] File .env dikonfigurasi dengan JWT_SECRET yang kuat
[ ] Docker image di-build
[ ] Migration dijalankan
[ ] Superadmin user dibuat (go run ./cmd/seed)
[ ] Aplikasi berjalan (curl /health returns OK)
[ ] Nginx dikonfigurasi sebagai reverse proxy
[ ] SSL/HTTPS aktif
[ ] CORS_ORIGINS diisi dengan domain production
[ ] Backup otomatis dijadwalkan
[ ] Fail2ban aktif
[ ] Security headers dikonfigurasi
[ ] Gzip compression aktif
[ ] Monitoring aktif (opsional)
[ ] Pengaturan SaaS (nama, logo, dll) dikonfigurasi via admin panel

=============================================
Catatan: Panduan ini dapat di-edit oleh superadmin
melalui panel Admin → Documentation jika ada
penyesuaian lebih lanjut sesuai kebutuhan server.
=============================================
', updated_at = NOW()
WHERE title = 'Panduan Deployment ke Server / VPS';
