# TaskMaster REST API

Platform manajemen tugas sederhana yang mengandalkan kapabilitas *Cloud-Native*, dibangun menggunakan arsitektur monolith ringan, *database* tanpa server (Neon.tech PostgreSQL), dan fitur kontainer secara otomatis (Docker + GitHub Actions) yang menghasilkan image super tipis (*lite*) ~15-20 MB.

## 🚀 Fitur Pendekatan Cloud-Native
- **Performa Tinggi**: Menggunakan **Golang (Gin)** yang terkenal dengan performa *routing* terbaik.
- **Micro-Container**: Dibangun menggunakan *Multi-Stage Build* langsung ke `alpine:latest`.
- **Integrasi Database Otomatis**: **GORM** dan koneksi *Serverless* via Neon.tech memberikan *developer-experience* luar biasa.
- **Enterprise-Ready**: Mendukung soft-deletes dan keamanan menggunakan `JWT`.

## 🛠 Cara Setup (Local Development)

### 1. Prasyarat
- Go 1.22+
- Docker Desktop (Opsional)
- PostgreSQL / *Account* Neon.tech

### 2. Konfigurasi
Lakukan pengaturan konfigurasi pada berkas statik:
```bash
cp .env.example .env
```
Isi konfigurasi pada file `.env` dengan kredensial PostgreSQL dan secret key Anda.

### 3. Menjalankan Aplikasi
Buka *port default* (8080) pada komputasi lokal dan panggil perintah berikut:
```bash
go mod tidy
go run cmd/api/main.go
# Atau menggunakan container Docker:
# docker build -t taskmaster-api .
# docker run -p 8080:8080 --env-file .env taskmaster-api
```

## 🧪 Strategi Testing
Proyek ini mengimplementasikan pengujian ganda (API testing / Unit Testing). Untuk memeriksa validasi sistem dan coverage:
```bash
go test -v ./...
```

## ⚙️ Skema Integrasi Berkelanjutan (CI/CD)
Secara otomatis akan melakukan langkah-langkah di bawah pada *branch `main`*:
- Melakukan pemeriksaan test (`go test`).
- *Login* otomasi ke arsitektur Docker Hub.
- Menyusun & menaruh *file* Image yang telah dibentuk (`push`).

*(Jangan lupa menyediakan secrets di Github: `DOCKERHUB_USERNAME` dan `DOCKERHUB_TOKEN`!)*

## ☁️ Deployment ke Render.com (via Docker Hub & Webhook)
Proyek ini dikonfigurasi menggunakan pipeline tingkat lanjut di mana **GitHub Actions** berperan sebagai *builder*, dan **Render** berperan murni sebagai *runner*.

### Konfigurasi Render
1. Login ke Render *Dashboard* dan buat layanan baru: **New Web Service**.
2. Pilih opsi **Deploy an existing image from a registry**.
3. Masukkan Image URL Docker Hub Anda (contoh: `indragiri21/taskmaster-api:latest`) dan klik Next.
4. Pilih tier *Free*, lalu pada bagian **Environment Variables**, masukkan rahasia Anda (`JWT_SECRET` & `DATABASE_URL`).
5. Pada bagian **Auto-Deploy**, atur ke **No**.
6. Simpan konfigurasi, salin URL dari bagian **Deploy Hook** di halaman *Settings* Render Anda.

### Konfigurasi Automasi (GitHub Actions)
1. Buka Repositori GitHub Anda -> **Settings** -> **Secrets and variables** -> **Actions**.
2. Tambahkan Secret baru bernama `RENDER_DEPLOY_HOOK` dan isi dengan URL Deploy Hook yang Anda salin dari Render.
3. Mulai sekarang, setiap kali Anda mem-*push* kode ke branch `main`, GitHub Actions akan:
   - Menjalankan unit tests (`go test`).
   - Melakukan *Build* dan *Push* layer Docker ke akun Docker Hub Anda.
   - Mengirim trigger (via webhook `curl`) ke Render.
   - Render secara instan akan menukar kontainernya ke versi image terbaru dari Docker Hub tanpa hambatan (*Zero-Downtime Deployment*).
