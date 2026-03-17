# Database Migration Guide

## Cara Mengubah/Menghapus Tabel yang Sudah Dieksekusi

### 1. Drop Tabel yang Ada (Manual)

Jika Anda ingin menghapus tabel yang sudah ada:

```sql
-- Koneksi ke MySQL dan jalankan:
DROP TABLE IF EXISTS scores;
DROP TABLE IF EXISTS users;
```

### 2. Menggunakan GORM Migrator

GORM menyediakan migrator untuk mengelola perubahan schema:

```go
// Drop tabel
db.Migrator().DropTable(&models.User{})
db.Migrator().DropTable(&models.Score{})

// Drop kolom tertentu
db.Migrator().DropColumn(&models.User{}, "phone")

// Rename kolom
db.Migrator().RenameColumn(&models.User{}, "name", "full_name")

// Add kolom
db.Migrator().AddColumn(&models.User{}, "username")

// Check apakah tabel ada
if db.Migrator().HasTable(&models.User{}) {
    // tabel ada
}
```

### 3. Reset Database (Development Only)

Cara paling mudah untuk development:

```sql
-- Drop seluruh database
DROP DATABASE scoring_app_db;

-- Buat ulang database
CREATE DATABASE scoring_app_db;
```

Kemudian jalankan aplikasi lagi, GORM akan membuat tabel baru sesuai model.

### 4. Migration File (Recommended for Production)

Untuk production, sebaiknya gunakan migration file terpisah seperti:

- golang-migrate/migrate
- gorm-migrator
- sql-migration files

## Best Practices

### Development

- Drop dan recreate database sesuka hati
- Gunakan auto migrate dari GORM

### Production

- **JANGAN** drop tabel production!
- Gunakan migration files dengan versioning
- Backup database sebelum migrate
- Test migration di staging dulu
- Gunakan transaction untuk migration

## Reset Database untuk Project Baru

Karena kita mengubah total struktur aplikasi, cara termudah:

1. **Stop aplikasi** jika sedang berjalan

2. **Drop database:**

   ```sql
   DROP DATABASE IF EXISTS scoring_app_db;
   CREATE DATABASE scoring_app_db;
   ```

3. **Jalankan aplikasi lagi** dengan model baru
   ```bash
   go run main.go
   ```

GORM akan otomatis membuat tabel baru sesuai model yang didefinisikan.

## Migrasi Bertahap (Jika Ada Data Production)

Jika Anda memiliki data production yang perlu dimigrasi:

```go
// 1. Buat tabel baru dulu
db.AutoMigrate(&models.NewModel{})

// 2. Migrasi data dari tabel lama ke baru
// ... kode migrasi data

// 3. Drop tabel lama
db.Migrator().DropTable(&models.OldModel{})
```

## Tools yang Bisa Digunakan

1. **golang-migrate/migrate** - Migration tool dengan version control
2. **GORM Auto Migration** - Sederhana untuk development
3. **SQL Migration Files** - Manual tapi fleksibel
4. **Liquibase/Flyway** - Enterprise-grade migration tools

Untuk project scoring_app ini, kita gunakan GORM Auto Migration karena masih development.
