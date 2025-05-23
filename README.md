# 🚀 PT_XYZ_Test – Studi Kasus Golang Developer (KREDIT PLUS)

Selamat datang di proyek studi kasus PT XYZ!  
Aplikasi ini merupakan implementasi backend service berbasis **Golang** yang dibangun dengan pendekatan **GitFlow** serta mengikuti praktik pengembangan modern.

---

## 📁 Struktur dan Fitur Utama

- ✅ Arsitektur modular menggunakan **repository pattern** dan **service layer**
- ✅ Menggunakan **Fiber (Go Web Framework)** untuk HTTP handling
- ✅ Koneksi ke database **MySQL** via **GORM**
- ✅ Dukungan upload file dengan validasi
- ✅ Dockerized untuk kemudahan deploy dan environment parity
- ✅ Mengadopsi **GitFlow Workflow** (`main`, `develop`, `feature/*`, `hotfix/*`)

---

## 📦 Kebutuhan Sistem

Sebelum menjalankan aplikasi, pastikan Anda memiliki:

- Docker & Docker Compose
- Git
- Go (hanya jika ingin menjalankan lokal tanpa Docker)

---

## 📚 File Tambahan

Semua file pendukung disimpan di dalam folder `extras/`:

| Jenis                     | File / Lokasi                    |
|--------------------------|----------------------------------|
| 💾 SQL Dump              | `extras/db.sql`                  |
| 🏗️ Gambar Arsitektur App | `extras/architecture.png`        |
| 🧩 Entity Relationship    | `extras/er-diagram.png`          |

---

## 🚀 Cara Menjalankan Proyek

### 🐳 Menggunakan Docker (Disarankan)

```bash
# 1. Clone repositori
git clone https://github.com/namamu/PT_XYZ_Test.git
cd PT_XYZ_Test

# 2. Jalankan Docker
docker-compose up --build
