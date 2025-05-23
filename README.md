# 🚀 PT_XYZ_Test – Studi Kasus Golang Developer (KREDIT PLUS)

Selamat datang di proyek studi kasus PT XYZ!  

---

## 📁 Struktur dan Fitur Utama

- ✅ Arsitektur modular menggunakan **repository pattern** dan **service layer**
- ✅ Menggunakan **Fiber (Go Web Framework)** untuk HTTP handling
- ✅ Koneksi ke database **MySQL** via **GORM**
- ✅ Dukungan upload file dengan validasi
- ✅ Dockerized
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
| 💾 SQL Dump              | `extras/test_pt_xyz.sql`                  |
| 🏗️ Gambar Arsitektur App | `extras/Arsitektur aplikasi.png`        |
| 🧩 Entity Relationship    | `extras/Struktur-Database.png`          |

---

## 🚀 Cara Menjalankan Proyek

### 🐳 Menggunakan Docker (Disarankan)

```bash
# 1. Clone repositori
git clone https://github.com/saldriAka/PT_XYZ_Test.git
cd PT_XYZ_Test

# 2. Jalankan Docker
docker-compose up --build
