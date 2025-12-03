# UAS PROGRES - Achievement Management System

Sistem manajemen prestasi mahasiswa menggunakan Go (Fiber), PostgreSQL, dan MongoDB.

## 📊 Database Schema

### PostgreSQL (7 Tabel)
1. `users` - Data pengguna (admin, dosen, mahasiswa)
2. `roles` - Role/peran pengguna
3. `permissions` - Hak akses sistem
4. `role_permissions` - Relasi many-to-many role dan permission
5. `students` - Data mahasiswa
6. `lecturers` - Data dosen
7. `achievement_references` - Referensi prestasi ke MongoDB

### MongoDB (1 Collection)
- `achievements` - Data detail prestasi mahasiswa (flexible schema)

## 🚀 Setup & Installation

### Prerequisites
- Go 1.21+
- PostgreSQL 14+
- MongoDB 6+

### Environment Variables
Buat file `.env`:
```env
DB_DSN=postgres://user:password@localhost:5432/dbname?sslmode=disable
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=achievements_db
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h
```

### Database Setup
```bash
# PostgreSQL - Buat tabel
psql -U your_user -d your_database -f migrations/000_create_tables.sql

# PostgreSQL - Insert sample data
psql -U your_user -d your_database -f migrations/002_insert_sample_data.sql
```

### Run Application
```bash
# Install dependencies
go mod tidy

# Build
go build -o app main.go

# Run
./app
# atau
go run main.go
```

Server akan berjalan di `http://localhost:4000`

## 📁 Project Structure
```
.
├── Domain/
│   ├── config/          # Database & JWT config
│   ├── middleware/      # Auth & role middleware
│   ├── model/
│   │   ├── Postgresql/  # PostgreSQL models
│   │   └── mongoDB/     # MongoDB models
│   ├── repository/      # Database operations
│   ├── route/           # API routes
│   └── service/         # Business logic
├── migrations/          # SQL migrations
├── main.go
└── .env
```

## 🔐 Authentication

### Login
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "mahasiswa@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "message": "Login berhasil",
  "token": "eyJhbGc...",
  "user": {
    "id": "uuid",
    "username": "mahasiswa1",
    "full_name": "Mahasiswa Satu",
    "email": "mahasiswa@example.com",
    "role_id": "uuid"
  }
}
```

### Logout
```bash
POST /auth/logout
Authorization: Bearer <token>
```

## 👥 Default Users

| Username   | Email                  | Password    | Role      |
|------------|------------------------|-------------|-----------|
| admin      | admin@example.com      | password123 | Admin     |
| dosen1     | dosen@example.com      | password123 | Dosen     |
| mahasiswa1 | mahasiswa@example.com  | password123 | Mahasiswa |

## 📝 API Documentation

Lihat file `DATABASE_SCHEMA.md` untuk detail skema database lengkap.

## ⚠️ Important Notes

1. **JANGAN** menambahkan tabel baru di PostgreSQL (hanya 7 tabel yang diizinkan)
2. **JANGAN** mengubah struktur tabel yang sudah ada
3. Token blacklist menggunakan in-memory storage (untuk production gunakan Redis)
4. MongoDB collection akan dibuat otomatis saat insert pertama

## 🔧 Development

### Build
```bash
go build -o app main.go
```

### Test Connection
```bash
# Test PostgreSQL
psql -U your_user -d your_database -c "SELECT version();"

# Test MongoDB
mongosh --eval "db.version()"
```

## 📚 Documentation Files

- `DATABASE_SCHEMA.md` - Skema database lengkap
- `PERBAIKAN.md` - Detail perbaikan yang dilakukan
- `RINGKASAN_PERBAIKAN.md` - Ringkasan perbaikan

## 🎯 Features

- ✅ Authentication & Authorization (JWT)
- ✅ Role-based Access Control (RBAC)
- ✅ Token Blacklist (Logout)
- ✅ User Management
- ✅ Achievement Management (PostgreSQL + MongoDB)
- ✅ Flexible Achievement Schema (MongoDB)

## 📄 License

UAS Project - Backend Lanjutan