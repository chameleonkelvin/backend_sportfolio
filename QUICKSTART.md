# Quick Start Guide

## Setup dan Jalankan Aplikasi

### 1. Konfigurasi Environment

Edit file `.env` dan sesuaikan dengan konfigurasi database MySQL Anda:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password_here
DB_NAME=scoring_app_db
```

### 2. Buat Database

Buat database MySQL dengan nama yang sesuai dengan konfigurasi di `.env`:

```sql
CREATE DATABASE scoring_app_db;
```

### 3. Install Dependencies

Dependencies sudah terinstall. Jika perlu update, jalankan:

```bash
go mod tidy
```

### 4. Jalankan Aplikasi

```bash
go run main.go
```

Atau gunakan Makefile:

```bash
make run
```

Aplikasi akan berjalan di `http://localhost:8080`

## Testing API dengan cURL

### 1. Health Check

```bash
curl http://localhost:8080/health
```

### 2. Create User

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Budi Santoso\",\"email\":\"budi@example.com\",\"phone\":\"081234567890\"}"
```

Response:

```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "name": "Budi Santoso",
    "email": "budi@example.com",
    "phone": "081234567890",
    "created_at": "2026-03-02T10:00:00Z",
    "updated_at": "2026-03-02T10:00:00Z"
  }
}
```

### 3. Get All Users

```bash
curl "http://localhost:8080/api/v1/users?page=1&page_size=10"
```

### 4. Get User by ID

```bash
curl http://localhost:8080/api/v1/users/1
```

### 5. Create Score

```bash
curl -X POST http://localhost:8080/api/v1/scores \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":1,\"category\":\"Matematika\",\"value\":85,\"max_value\":100,\"description\":\"Ujian Tengah Semester\"}"
```

Response:

```json
{
  "success": true,
  "message": "Score created successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "category": "Matematika",
    "value": 85,
    "max_value": 100,
    "description": "Ujian Tengah Semester",
    "created_at": "2026-03-02T10:05:00Z",
    "updated_at": "2026-03-02T10:05:00Z"
  }
}
```

### 6. Create Multiple Scores

```bash
# Score untuk Fisika
curl -X POST http://localhost:8080/api/v1/scores \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":1,\"category\":\"Fisika\",\"value\":90,\"max_value\":100,\"description\":\"Ujian Akhir Semester\"}"

# Score untuk Kimia
curl -X POST http://localhost:8080/api/v1/scores \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":1,\"category\":\"Kimia\",\"value\":78,\"max_value\":100,\"description\":\"Praktikum\"}"
```

### 7. Get All Scores

```bash
curl "http://localhost:8080/api/v1/scores?page=1&page_size=10"
```

### 8. Get Scores by User ID

```bash
curl http://localhost:8080/api/v1/scores/user/1
```

### 9. Get Average Score for User

```bash
curl http://localhost:8080/api/v1/scores/user/1/average
```

Response:

```json
{
  "success": true,
  "message": "Average score calculated successfully",
  "data": {
    "user_id": 1,
    "average": 84.33
  }
}
```

### 10. Update Score

```bash
curl -X PUT http://localhost:8080/api/v1/scores/1 \
  -H "Content-Type: application/json" \
  -d "{\"category\":\"Matematika\",\"value\":95,\"max_value\":100,\"description\":\"Ujian Tengah Semester - Revisi\"}"
```

### 11. Update User

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Budi Santoso Update\",\"email\":\"budi.new@example.com\",\"phone\":\"089876543210\"}"
```

### 12. Delete Score

```bash
curl -X DELETE http://localhost:8080/api/v1/scores/1
```

### 13. Delete User

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Testing dengan Postman

Import collection berikut ke Postman:

### Base URL

```
http://localhost:8080
```

### Headers

```
Content-Type: application/json
```

### Endpoints Testing Sequence

1. **Health Check** - GET `/health`
2. **Create User** - POST `/api/v1/users`
3. **Get All Users** - GET `/api/v1/users?page=1&page_size=10`
4. **Get User by ID** - GET `/api/v1/users/:id`
5. **Create Score** - POST `/api/v1/scores`
6. **Get All Scores** - GET `/api/v1/scores?page=1&page_size=10`
7. **Get User Scores** - GET `/api/v1/scores/user/:user_id`
8. **Get Average Score** - GET `/api/v1/scores/user/:user_id/average`
9. **Update Score** - PUT `/api/v1/scores/:id`
10. **Update User** - PUT `/api/v1/users/:id`
11. **Delete Score** - DELETE `/api/v1/scores/:id`
12. **Delete User** - DELETE `/api/v1/users/:id`

## Troubleshooting

### Database Connection Error

Pastikan:

- MySQL sudah berjalan
- Database sudah dibuat
- Kredensial di `.env` sudah benar
- Host dan port sesuai

### Port Already in Use

Ubah `SERVER_PORT` di file `.env`:

```env
SERVER_PORT=8081
```

### Permission Denied

Jika di Linux/Mac, pastikan file memiliki permission yang benar:

```bash
chmod +x scoring_app
```
