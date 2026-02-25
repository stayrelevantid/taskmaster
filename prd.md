# PRD: TaskMaster API (Cloud-Native Edition)

## 1. Identitas Proyek
- **Nama Proyek:** TaskMaster REST API
- **Tujuan:** Platform manajemen tugas (CRUD) yang aman, tervalidasi, dan siap dideploy di lingkungan kontainer (Docker/Kubernetes).
- **Status:** Production-Ready (Portfolio Prototype).

## 2. Tech Stack & Requirements

| Komponen | Teknologi | Keterangan |
| --- | --- | --- |
| **Language** | Go 1.22+ | Standard industri untuk performa tinggi. |
| **Framework** | Gin Gonic | Routing HTTP yang cepat dan ringan. |
| **Database** | PostgreSQL | Managed via Neon.tech (Serverless). |
| **ORM** | GORM | Memudahkan manajemen skema dan relasi. |
| **Auth** | JWT | Keamanan stateless (kunci rahasia via Environment Variable). |
| **Container** | Docker | Multi-stage build (Ukuran image ~15-20 MB). |
| **CI/CD** | GitHub Actions | Otomatisasi Build, Test, dan Push ke Registry. |

## 3. Fitur Utama & Spesifikasi Teknis

### A. Keamanan & Validasi (Enterprise Standard)
- **Middleware Auth:** Memastikan akses ke resource sensitif hanya dilakukan oleh user yang memiliki token valid.
- **Input Sanitization:** Menggunakan struct tagging (`binding:"required,min=3"`) untuk mencegah data korup.
- **Standard Response:** Format JSON seragam (`success`, `message`, `data`) untuk memudahkan integrasi frontend.

### B. Matrix API Endpoint

| Method | Path | Auth | Deskripsi |
| --- | --- | --- | --- |
| GET | `/health` | No | Monitoring kesehatan aplikasi. |
| POST | `/login` | No | Mendapatkan JWT (Dummy Credentials). |
| GET | `/api/v1/tasks` | Yes | Mengambil semua daftar tugas. |
| POST | `/api/v1/tasks` | Yes | Membuat tugas baru. |
| PUT | `/api/v1/tasks/:id` | Yes | Memperbarui judul/status tugas. |
| DELETE | `/api/v1/tasks/:id` | Yes | Menghapus tugas (Soft Delete). |

## 4. Fase Eksekusi (Execution Flow)

- **Fase 1: Inisialisasi Lokal**
  - Setup modul Go: `go mod init`.
  - Instalasi dependensi (Gin, GORM, JWT).
  - Setup database di Neon.tech dan ambil koneksi string (DSN).
- **Fase 2: Development & Kode**
  - Implementasi `models/task.go` dengan validasi.
  - Implementasi `middleware/auth.go` untuk proteksi route.
  - Implementasi `main.go` yang menggabungkan seluruh logic.
- **Fase 3: Kontainerisasi (Docker)**
  - Membuat `.dockerignore` untuk mengecualikan file sampah.
  - Membuat `Dockerfile` multi-stage:
    - *Stage 1 (Builder):* Compile binary.
    - *Stage 2 (Runtime):* Menjalankan binary di image Alpine yang super kecil.
- **Fase 4: CI/CD & Deployment**
  - Setup GitHub Secrets (`DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN`).
  - Konfigurasi GitHub Actions untuk automated deployment.
  - Deploy image ke Koyeb/Render dan injeksi Environment Variables.
- **Fase 5: Unit Testing & Dokumentasi**
  - Refaktorisasi main.go untuk *testability*.
  - Pembuatan automated testing `go test`.
  - Pembuatan `README.md` komprehensif.

## 5. Environment Variables (Secret Management)
Untuk alasan keamanan, aplikasi tidak boleh berjalan tanpa variabel berikut:
- `DATABASE_URL`: Alamat koneksi ke PostgreSQL.
- `JWT_SECRET`: Kunci privat untuk enkripsi token.
- `PORT`: Port aplikasi (default: 8080).

## 6. Prosedur Pengujian (Testing Strategy)
- **Unit Testing:** Menjalankan `go test ./...` untuk memastikan fungsi logika benar.
- **Integration Testing:** Menggunakan Postman Collection untuk menguji alur Login -> Ambil Token -> CRUD Task.
- **Deployment Testing:** Verifikasi endpoint `/health` setelah aplikasi live di cloud.

## 7. Instruksi Replikasi Cepat (CLI)
Cukup jalankan satu blok perintah ini untuk inisialisasi awal:

```bash
# Buat folder dan inisialisasi
mkdir taskmaster-api && cd taskmaster-api
go mod init github.com/yourusername/taskmaster-api

# Install dependensi wajib
go get -u github.com/gin-gonic/gin gorm.io/gorm gorm.io/driver/postgres github.com/golang-jwt/jwt/v5

# Buat struktur folder
mkdir -p cmd/api internal/middleware internal/models .github/workflows
```