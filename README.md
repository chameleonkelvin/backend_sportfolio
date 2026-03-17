# Match Management System API

API untuk aplikasi manajemen pertandingan olahraga beregu (badminton, tenis, dll) menggunakan Golang dengan arsitektur yang clean dan terstruktur.

## 🎯 Overview

Aplikasi ini adalah sistem manajemen pertandingan yang memungkinkan:

- **Event Organizer** membuat dan mengelola event pertandingan
- **Admin** mengelola sistem dan user
- **Player** berpartisipasi dalam pertandingan
- Sistem scoring otomatis untuk pertandingan beregu (2v2)
- Generate ronde pertandingan dengan multiple lapangan

## 🛠️ Tech Stack

- **Go** (Golang) - Backend language
- **Gin Gonic** - Web Framework
- **GORM** - ORM
- **MySQL** - Database
- **Godotenv** - Environment Configuration

## 📁 Project Structure

```
scoring_app/
├── config/                      # Konfigurasi aplikasi
├── controllers/                 # HTTP handlers (TODO: perlu dibuat)
├── database/                    # Database connection & migrations
│   ├── database.go             # Connection setup
│   ├── migration.go            # Auto migration & reset
│   └── seeder.go               # Seed initial data
├── models/                      # Database models
│   ├── account_type.go         # Model untuk jenis akun
│   ├── user.go                 # Model user dengan auth
│   ├── match_event.go          # Model event pertandingan
│   ├── match_player.go         # Model pemain
│   └── match_round.go          # Model ronde pertandingan
├── repositories/                # Data access layer
│   ├── account_type_repository.go
│   ├── user_repository.go
│   ├── match_event_repository.go
│   ├── match_player_repository.go
│   └── match_round_repository.go
├── routes/                      # Route definitions (TODO: perlu diupdate)
├── services/                    # Business logic layer (TODO: perlu dibuat)
├── validators/                  # Request validation & responses
├── .env                        # Environment variables
├── .env.example                # Template environment
├── main.go                     # Application entry point
├── DATABASE_SCHEMA.md          # Detail schema database
├── MIGRATION_GUIDE.md          # Panduan migrasi
└── SETUP_NEW_DATABASE.md       # Panduan setup database baru
```

## 📊 Database Schema

### Tables:

1. **account_types** - Jenis akun (admin, organizer, player)
2. **users** - User dengan authentication
3. **match_events** - Event pertandingan
4. **match_players** - Pemain dalam event
5. **match_rounds** - Ronde pertandingan (siapa vs siapa)

Lihat detail di [DATABASE_SCHEMA.md](DATABASE_SCHEMA.md)

## 🚀 Setup & Installation

### Prerequisites

- Go 1.21 atau lebih tinggi
- MySQL 8.0 atau lebih tinggi

### Step-by-Step Installation

1. **Clone repository ini**

2. **Copy dan edit file `.env`:**

```bash
cp .env.example .env
```

Edit `.env` dan sesuaikan dengan konfigurasi MySQL Anda:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=scoring_app_db
```

3. **Buat database MySQL:**

```sql
CREATE DATABASE scoring_app_db;
```

4. Install dependencies:

```bash
go mod tidy
```

4. **Reset database dan seed data:**

```bash
# Reset database (drop & recreate tables) + seed initial data
go run main.go -reset

# Server akan running di http://localhost:8080
```

## 🎮 Command Line Options

### Reset Database (Development)

```bash
go run main.go -reset
```

Drop semua tabel dan buat ulang + seed data account types.

### Seed Data Saja

```bash
go run main.go -seed
```

Hanya menambahkan data seed tanpa drop tabel.

### Normal Run

```bash
go run main.go
```

Jalankan aplikasi dengan auto migrate (tidak drop data).

## 📡 API Endpoints (TODO)

> **Note:** Controllers dan services masih perlu diimplementasikan.
> Saat ini hanya tersedia health check endpoint.

### Current Endpoints

#### Health Check

- `GET /health` - Check API status

### Planned Endpoints

#### Authentication

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login (JWT)
- `POST /api/v1/auth/forgot-password` - Request OTP
- `POST /api/v1/auth/reset-password` - Reset password with OTP

#### Account Types

- `GET /api/v1/account-types` - Get all account types

#### Users

- `GET /api/v1/users` - Get all users (pagination)
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user profile
- `DELETE /api/v1/users/:id` - Delete user (soft delete)

#### Match Events

- `POST /api/v1/match-events` - Create new event
- `GET /api/v1/match-events` - Get all events (pagination)
- `GET /api/v1/match-events/:id` - Get event detail with players & rounds
- `PUT /api/v1/match-events/:id` - Update event
- `DELETE /api/v1/match-events/:id` - Delete event
- `GET /api/v1/match-events/upcoming` - Get upcoming events
- `GET /api/v1/match-events/my-events` - Get events created by current user

#### Match Players

- `POST /api/v1/match-players` - Add player to event
- `POST /api/v1/match-players/batch` - Add multiple players
- `GET /api/v1/match-players/event/:match_id` - Get all players in event
- `PUT /api/v1/match-players/:id` - Update player
- `DELETE /api/v1/match-players/:id` - Remove player

#### Match Rounds

- `POST /api/v1/match-rounds` - Create round
- `POST /api/v1/match-rounds/generate` - Auto-generate rounds for event
- `GET /api/v1/match-rounds/event/:match_id` - Get all rounds for event
- `GET /api/v1/match-rounds/:id` - Get round detail with players
- `PUT /api/v1/match-rounds/:id` - Update round scores
- `DELETE /api/v1/match-rounds/:id` - Delete round

## 📝 Request Examples (Planned)

### Health Check

```bash
curl http://localhost:8080/health
```

Response:

```json
{
  "status": "ok",
  "message": "Scoring App API is running",
  "models": "account_types, users, match_events, match_players, match_rounds"
}
```

### Register User (Organizer)

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "account_type_id": "organizer",
    "username": "ahmad_badminton",
    "full_name": "Ahmad Santoso",
    "email": "ahmad@example.com",
    "password": "SecurePassword123!",
    "birth_date": "1990-05-15"
  }'
```

### Create Match Event

```bash
curl -X POST http://localhost:8080/api/v1/match-events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt_token>" \
  -d '{
    "name": "Turnamen Bulutangkis Ramadan 2026",
    "total_courts": 3,
    "game_type": "Badminton",
    "location": "GOR Pancasila Jakarta",
    "play_date": "2026-04-15T09:00:00Z",
    "total_players": 16,
    "team_type": "men_double"
  }'
```

### Add Players to Event

```bash
curl -X POST http://localhost:8080/api/v1/match-players/batch \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt_token>" \
  -d '{
    "match_id": "event-uuid-123",
    "players": [
      {"name": "Ahmad", "gender": 1},
      {"name": "Budi", "gender": 1},
      {"name": "Candra", "gender": 1},
      {"name": "Deni", "gender": 1}
    ]
  }'
```

### Create Match Round

```bash
curl -X POST http://localhost:8080/api/v1/match-rounds \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt_token>" \
  -d '{
    "match_id": "event-uuid-123",
    "round_number": 1,
    "court": 1,
    "team_a_player_1_id": 1,
    "team_a_player_2_id": 2,
    "team_b_player_1_id": 3,
    "team_b_player_2_id": 4,
    "score_a": 21,
    "score_b": 15
  }'
```

### Get Event with Details

```bash
curl http://localhost:8080/api/v1/match-events/event-uuid-123 \
  -H "Authorization: Bearer <jwt_token>"
```

## 📤 Response Format

### Success Response

```json
{
  "success": true,
  "message": "Success message",
  "data": {}
}
```

### Error Response

```json
{
  "success": false,
  "message": "Error message",
  "error": "Error details"
}
```

### Pagination Response

```json
{
  "success": true,
  "message": "Success message",
  "data": [],
  "meta": {
    "page": 1,
    "page_size": 10,
    "total": 100,
    "total_page": 10
  }
}
```

## 🗄️ Database Models

### Account Types

- `admin` - Full system access
- `organizer` - Create and manage events
- `player` - Participate in matches

### Users

Authentication dengan bcrypt password hashing, support OTP untuk reset password.

### Match Events

Event pertandingan dengan detail lapangan, lokasi, dan tanggal.

### Match Players

Daftar pemain dalam event dengan tracking skor sementara.

### Match Rounds

Detail ronde pertandingan: siapa vs siapa, di lapangan berapa, skor berapa.

**Lihat detail lengkap:** [DATABASE_SCHEMA.md](DATABASE_SCHEMA.md)

## 🏗️ Implementasi Selanjutnya

Aplikasi ini sudah memiliki:

- ✅ Database models
- ✅ Repository layer (data access)
- ✅ Migration & seeding system
- ✅ Project structure

Yang perlu diimplementasikan:

- ⏳ **Services** - Business logic untuk setiap module
- ⏳ **Controllers** - HTTP handlers untuk API endpoints
- ⏳ **Validators** - Request validation
- ⏳ **Authentication** - JWT middleware
- ⏳ **Routes** - Setup API routes
- ⏳ **Password Hashing** - bcrypt implementation
- ⏳ **UUID/ULID** - Generate unique IDs

### Recommended Libraries untuk Implementasi

```bash
# JWT Authentication
go get github.com/golang-jwt/jwt/v5

# Password Hashing
go get golang.org/x/crypto/bcrypt

# UUID Generation
go get github.com/google/uuid

# atau ULID (sortable)
go get github.com/oklog/ulid/v2
```

## 💻 Development

### Run in Development Mode

```bash
go run main.go
```

### Reset Database (Development Only!)

```bash
go run main.go -reset
```

### Build for Production

```bash
go build -o scoring_app.exe
```

### Build with Optimization

```bash
go build -ldflags="-s -w" -o scoring_app.exe
```

## 📚 Documentation

- [README.md](README.md) - Overview & quick start
- [DATABASE_SCHEMA.md](DATABASE_SCHEMA.md) - Detailed database schema & ERD
- [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) - Database migration guide
- [SETUP_NEW_DATABASE.md](SETUP_NEW_DATABASE.md) - Setup guide untuk database baru
- [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - Project structure explanation

## ⚠️ Important Notes

1. **ID Generation**: Models menggunakan `varchar(50)` untuk ID. Gunakan UUID atau ULID untuk generate ID unik.
2. **Authentication**: Password harus di-hash menggunakan bcrypt sebelum disimpan.
3. **Soft Delete**: Hampir semua tabel support soft delete (kecuali match_players dan match_rounds).
4. **Migration**: Gunakan `-reset` flag HANYA di development. Production harus menggunakan migration tools yang proper.

## 🤝 Contributing

Contributions are welcome! Please follow:

1. Fork repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## 📄 License

MIT
