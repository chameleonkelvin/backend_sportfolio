# Database Schema (ERD)

## Entity Relationship Diagram

```
┌─────────────────┐
│ account_types   │
├─────────────────┤
│ id (PK)         │ varchar(50)
│ name            │ varchar(100)
│ description     │ text
│ created_at      │ datetime
│ updated_at      │ datetime
│ deleted_at      │ datetime
└─────────────────┘
         │
         │ 1:N
         ▼
┌─────────────────┐
│ users           │
├─────────────────┤
│ id (PK)         │ varchar(50)
│ account_type_id │ varchar(50) [FK]
│ username        │ varchar(100) [UNIQUE]
│ full_name       │ varchar(255)
│ email           │ varchar(255) [UNIQUE]
│ password_hash   │ varchar(255)
│ otp_code        │ varchar(10)
│ birth_date      │ date
│ created_at      │ datetime
│ updated_at      │ datetime
│ deleted_at      │ datetime
└─────────────────┘
         │
         │ 1:N (pembuat event)
         ▼
┌─────────────────┐
│ match_events    │
├─────────────────┤
│ id (PK)         │ varchar(50)
│ user_id (FK)    │ varchar(50)
│ name            │ varchar(255)
│ total_courts    │ int
│ game_type       │ varchar(100)
│ location        │ varchar(255)
│ play_date       │ datetime
│ total_players   │ int
│ team_type       │ varchar(50)
│ created_at      │ datetime
│ updated_at      │ datetime
│ deleted_at      │ datetime
└─────────────────┘
         │
         │ 1:N
         ▼
┌─────────────────┐        ┌─────────────────┐
│ match_players   │        │ match_rounds    │
├─────────────────┤        ├─────────────────┤
│ id (PK)         │◄───────│ match_id (FK)   │
│ match_id (FK)   │        │ round_number    │ int
│ name            │        │ court           │ int
│ gender          │ int    │ team_a_player_1 │ int [FK]
│ temp_score      │ int    │ team_a_player_2 │ int [FK]
│ created_at      │        │ score_a         │ int
└─────────────────┘        │ team_b_player_1 │ int [FK]
                           │ team_b_player_2 │ int [FK]
                           │ score_b         │ int
                           │ created_at      │
                           └─────────────────┘
```

## Tables Description

### account_types

Menyimpan jenis-jenis akun yang ada di sistem.

**Data Awal (Seed):**

- `admin` - Administrator (Full access)
- `organizer` - Event Organizer (Buat dan kelola event)
- `player` - Player (Ikut dalam pertandingan)

### users

Menyimpan data user yang menggunakan aplikasi.

**Fields:**

- `id` - Unique identifier (UUID/ULID recommended)
- `account_type_id` - Jenis akun user
- `username` - Username unik untuk login
- `full_name` - Nama lengkap user
- `email` - Email unik
- `password_hash` - Password dalam bentuk hash (bcrypt)
- `otp_code` - Kode OTP untuk reset password/2FA
- `birth_date` - Tanggal lahir

**Relations:**

- Belongs to: `account_types`
- Has many: `match_events` (sebagai pembuat event)

### match_events

Menyimpan data event pertandingan yang dibuat oleh organizer.

**Fields:**

- `id` - Unique identifier
- `user_id` - Pembuat event (organizer)
- `name` - Nama event (contoh: "Turnamen Bulutangkis Ramadan 2026")
- `total_courts` - Jumlah lapangan yang tersedia
- `game_type` - Jenis permainan (Badminton, Tenis, dll)
- `location` - Lokasi event
- `play_date` - Tanggal dan waktu main
- `total_players` - Total pemain yang ikut
- `team_type` - Tipe tim:
  - "men_double" - Ganda Putra
  - "women_double" - Ganda Putri
  - "mix_double" - Ganda Campuran

**Relations:**

- Belongs to: `users`
- Has many: `match_players`
- Has many: `match_rounds`

### match_players

Menyimpan data pemain yang ikut dalam suatu event.

**Fields:**

- `id` - Auto increment ID
- `match_id` - Event yang diikuti
- `name` - Nama pemain
- `gender` - Jenis kelamin (0 = perempuan, 1 = pria)
- `temp_score` - Skor sementara/akumulatif sepanjang event

**Relations:**

- Belongs to: `match_events`
- Referenced by: `match_rounds` (4x foreign key untuk team A & B)

### match_rounds

Menyimpan data ronde pertandingan (siapa vs siapa di lapangan berapa).

**Fields:**

- `id` - Auto increment ID
- `match_id` - Event yang dimainkan
- `round_number` - Nomor ronde (1, 2, 3, dst)
- `court` - Nomor lapangan
- `team_a_player_1_id` - Player 1 tim A
- `team_a_player_2_id` - Player 2 tim A
- `score_a` - Skor tim A
- `team_b_player_1_id` - Player 1 tim B
- `team_b_player_2_id` - Player 2 tim B
- `score_b` - Skor tim B

**Relations:**

- Belongs to: `match_events`
- References: `match_players` (4 foreign keys)

## Indexes

**Performance Indexes:**

```sql
-- users table
CREATE INDEX idx_users_account_type ON users(account_type_id);
CREATE UNIQUE INDEX idx_users_username ON users(username);
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- match_events table
CREATE INDEX idx_match_events_user ON match_events(user_id);
CREATE INDEX idx_match_events_play_date ON match_events(play_date);
CREATE INDEX idx_match_events_deleted_at ON match_events(deleted_at);

-- match_players table
CREATE INDEX idx_match_players_match ON match_players(match_id);

-- match_rounds table
CREATE INDEX idx_match_rounds_match ON match_rounds(match_id);
CREATE INDEX idx_match_rounds_players ON match_rounds(team_a_player_1_id, team_a_player_2_id, team_b_player_1_id, team_b_player_2_id);
```

## Data Flow Example

### Skenario: Membuat Event Badminton

1. **User Register & Login**
   - User register dengan account_type = "organizer"
   - User login dan dapat token

2. **Create Match Event**

   ```json
   {
     "name": "Turnamen Bulutangkis Ramadan 2026",
     "total_courts": 3,
     "game_type": "Badminton",
     "location": "GOR Pancasila",
     "play_date": "2026-04-15T09:00:00Z",
     "total_players": 16,
     "team_type": "men_double"
   }
   ```

3. **Add Players**

   ```json
   [
     {"name": "Ahmad", "gender": 1},
     {"name": "Budi", "gender": 1},
     {"name": "Candra", "gender": 1},
     ...
   ]
   ```

4. **Generate Rounds**
   - Sistem generate pasangan tim secara otomatis atau manual
   - Setiap ronde: 4 pemain (2v2) di suatu lapangan

5. **Input Scores**
   - Update skor untuk setiap ronde
   - Update temp_score pemain

6. **View Results**
   - Lihat hasil per ronde
   - Lihat ranking berdasarkan temp_score

## Business Rules

### User Management

1. Email dan username harus unique
2. Password disimpan dalam bentuk hash (bcrypt)
3. OTP code digunakan untuk reset password
4. Soft delete (deleted_at) untuk data retention

### Match Events

1. Hanya user dengan account_type "organizer" atau "admin" yang bisa create event
2. PlayDate harus di masa depan
3. Total players harus sesuai dengan jumlah pemain yang didaftarkan
4. Team type menentukan validasi gender pemain

### Match Rounds

1. Setiap ronde harus ada 4 pemain (2v2)
2. Player tidak boleh main melawan diri sendiri
3. Skor harus >= 0
4. Court number harus dalam range 1 - total_courts

### Scoring System

1. temp_score di match_players diupdate setiap selesai ronde
2. Pemenang ronde mendapat poin lebih tinggi
3. Ranking ditentukan dari total temp_score

## Migration Commands

### Reset Database (DEVELOPMENT ONLY!)

```bash
go run main.go -reset
```

### Seed Data Only

```bash
go run main.go -seed
```

### Normal Migration

```bash
go run main.go
```

## SQL Schema (Reference)

```sql
CREATE TABLE account_types (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME
);

CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,
    account_type_id VARCHAR(50) NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    otp_code VARCHAR(10),
    birth_date DATE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME,
    FOREIGN KEY (account_type_id) REFERENCES account_types(id),
    INDEX idx_account_type (account_type_id),
    INDEX idx_deleted (deleted_at)
);

CREATE TABLE match_events (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    total_courts INT NOT NULL,
    game_type VARCHAR(100) NOT NULL,
    location VARCHAR(255) NOT NULL,
    play_date DATETIME NOT NULL,
    total_players INT NOT NULL,
    team_type VARCHAR(50) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id),
    INDEX idx_user (user_id),
    INDEX idx_deleted (deleted_at)
);

CREATE TABLE match_players (
    id INT PRIMARY KEY AUTO_INCREMENT,
    match_id VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    gender INT NOT NULL COMMENT '0 = perempuan, 1 = pria',
    temp_score INT DEFAULT 0,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (match_id) REFERENCES match_events(id),
    INDEX idx_match (match_id)
);

CREATE TABLE match_rounds (
    id INT PRIMARY KEY AUTO_INCREMENT,
    match_id VARCHAR(50) NOT NULL,
    round_number INT NOT NULL,
    court INT NOT NULL,
    team_a_player_1_id INT NOT NULL,
    team_a_player_2_id INT NOT NULL,
    score_a INT DEFAULT 0,
    team_b_player_1_id INT NOT NULL,
    team_b_player_2_id INT NOT NULL,
    score_b INT DEFAULT 0,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (match_id) REFERENCES match_events(id),
    FOREIGN KEY (team_a_player_1_id) REFERENCES match_players(id),
    FOREIGN KEY (team_a_player_2_id) REFERENCES match_players(id),
    FOREIGN KEY (team_b_player_1_id) REFERENCES match_players(id),
    FOREIGN KEY (team_b_player_2_id) REFERENCES match_players(id),
    INDEX idx_match (match_id)
);
```
