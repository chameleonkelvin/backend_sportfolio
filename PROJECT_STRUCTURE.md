# Struktur Project Scoring App

```
scoring_app/
│
├── config/
│   └── config.go                 # Konfigurasi aplikasi dan environment
│
├── controllers/
│   ├── score_controller.go       # Handler untuk score endpoints
│   └── user_controller.go        # Handler untuk user endpoints
│
├── database/
│   ├── database.go               # Koneksi database
│   └── migration.go              # Auto migration models
│
├── models/
│   ├── score.go                  # Model Score
│   └── user.go                   # Model User
│
├── repositories/
│   ├── score_repository.go       # Data access layer untuk Score
│   └── user_repository.go        # Data access layer untuk User
│
├── routes/
│   └── routes.go                 # Definisi routing API
│
├── services/
│   ├── score_service.go          # Business logic untuk Score
│   └── user_service.go           # Business logic untuk User
│
├── validators/
│   ├── response.go               # Struktur response API
│   ├── score_validator.go        # Validasi request Score
│   └── user_validator.go         # Validasi request User
│
├── .env                          # Environment variables (jangan commit!)
├── .env.example                  # Template environment variables
├── .gitignore                    # File yang diabaikan Git
├── docker-compose.yml            # Docker Compose untuk MySQL
├── go.mod                        # Go module dependencies
├── go.sum                        # Go module checksums
├── main.go                       # Entry point aplikasi
├── Makefile                      # Make commands
├── QUICKSTART.md                 # Panduan quick start
└── README.md                     # Dokumentasi utama

```

## Deskripsi Folder

### config/

Berisi file konfigurasi aplikasi, termasuk:

- Membaca environment variables
- Konfigurasi database
- Konfigurasi server
- Helper functions untuk config

### controllers/

Layer yang menangani HTTP requests dan responses:

- Menerima request dari client
- Validasi input menggunakan validators
- Memanggil services untuk business logic
- Mengembalikan response ke client

### database/

Mengelola koneksi dan operasi database:

- Inisialisasi koneksi database
- Auto migration untuk membuat/update tabel
- Helper functions untuk database operations

### models/

Definisi struktur data (entities):

- Representasi tabel database
- Relasi antar tabel
- GORM tags untuk mapping

### repositories/

Data Access Layer (DAL):

- Interface untuk operasi database
- Implementasi CRUD operations
- Query kompleks untuk data retrieval
- Abstraksi dari database operations

### routes/

Definisi endpoint API:

- Grouping routes
- Mapping URL ke controllers
- Middleware setup (jika ada)

### services/

Business Logic Layer:

- Validasi business rules
- Orchestration antar repositories
- Kalkulasi dan transformasi data
- Error handling

### validators/

Request validation dan response structure:

- Struct untuk request validation
- Binding rules menggunakan Gin
- Standardisasi format response
- Pagination structure

## Flow Aplikasi

```
Client Request
      ↓
Routes (routes.go)
      ↓
Controller (controllers/)
      ↓
Validator (validators/)
      ↓
Service (services/)
      ↓
Repository (repositories/)
      ↓
Model (models/)
      ↓
Database (MySQL)
      ↓
Response back to Client
```

## Dependencies

### Main Dependencies

- **github.com/gin-gonic/gin** - Web framework
- **gorm.io/gorm** - ORM library
- **gorm.io/driver/mysql** - MySQL driver untuk GORM
- **github.com/joho/godotenv** - Load environment variables

### Supporting Dependencies

- **github.com/go-playground/validator/v10** - Input validation (included with Gin)
- **github.com/go-sql-driver/mysql** - MySQL driver (included with GORM MySQL)

## Best Practices yang Diterapkan

1. **Separation of Concerns**
   - Setiap layer memiliki tanggung jawab yang jelas
   - Tidak ada mixing logic antar layer

2. **Dependency Injection**
   - Service dan Repository di-inject ke controller
   - Memudahkan testing dan maintenance

3. **Interface-based Design**
   - Repository menggunakan interface
   - Service menggunakan interface
   - Memudahkan mocking untuk testing

4. **Error Handling**
   - Error handling di setiap layer
   - Standardisasi error response

5. **Configuration Management**
   - Environment-based configuration
   - Tidak ada hard-coded values

6. **Database Best Practices**
   - Soft delete menggunakan GORM
   - Indexing pada kolom yang sering di-query
   - Foreign key relationships

7. **API Design**
   - RESTful API conventions
   - Consistent response format
   - Pagination support
   - Proper HTTP status codes

## Cara Extend Aplikasi

### Menambah Model Baru

1. Buat file di `models/` (misal: `category.go`)
2. Tambahkan model ke `database/migration.go`
3. Buat repository di `repositories/`
4. Buat service di `services/`
5. Buat controller di `controllers/`
6. Buat validator di `validators/`
7. Tambahkan routes di `routes/routes.go`

### Menambah Endpoint Baru

1. Tambahkan method di controller yang sesuai
2. Tambahkan route di `routes/routes.go`
3. Implementasi logic di service jika perlu
4. Update dokumentasi di README.md

### Menambah Business Logic

1. Tambahkan method di service interface
2. Implementasi method di service struct
3. Gunakan method tersebut di controller
4. Test functionality

## Environment Variables

| Variable    | Description    | Default        |
| ----------- | -------------- | -------------- |
| DB_HOST     | MySQL host     | localhost      |
| DB_PORT     | MySQL port     | 3306           |
| DB_USER     | MySQL user     | root           |
| DB_PASSWORD | MySQL password | -              |
| DB_NAME     | Database name  | scoring_app_db |
| SERVER_PORT | Server port    | 8080           |
| SERVER_HOST | Server host    | localhost      |
| APP_ENV     | Environment    | development    |

## Testing

Untuk menambahkan testing:

1. Buat file `*_test.go` di setiap package
2. Gunakan `testing` package dari Go
3. Mock repositories menggunakan interface
4. Jalankan dengan `go test ./...`

## Production Deployment

1. Set `APP_ENV=production` di `.env`
2. Build dengan optimizations: `make build-prod`
3. Setup reverse proxy (nginx/apache)
4. Setup SSL certificate
5. Configure firewall
6. Setup logging dan monitoring
7. Database backup strategy
8. Load balancing jika perlu
