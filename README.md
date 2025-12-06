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

### Achievement Endpoints

#### FR-010: View All Achievements (Admin)
```bash
GET /api/v1/achievements?page=1&limit=10&status=verified&sort=created_at&order=desc
Authorization: Bearer <token>
Permission: read_achievements
```

Query Parameters:
- `page` - Halaman (default: 1)
- `limit` - Jumlah per halaman (default: 10, max: 100)
- `status` - Filter by status (draft, submitted, verified, rejected)
- `student_id` - Filter by student UUID
- `sort` - Sort by field (created_at, submitted_at, verified_at, updated_at)
- `order` - Sort order (asc, desc)

Response:
```json
{
  "message": "Berhasil mengambil data achievements",
  "data": {
    "achievements": [
      {
        "reference_id": "uuid",
        "achievement_id": "mongo_id",
        "student_id": "NIM123",
        "program_study": "Teknik Informatika",
        "status": "verified",
        "submitted_at": "2024-12-04T10:00:00Z",
        "verified_at": "2024-12-04T11:00:00Z",
        "verified_by": "lecturer-uuid",
        "rejection_note": null,
        "achievement": { ... },
        "created_at": "2024-12-03T09:00:00Z"
      }
    ],
    "pagination": {
      "total": 50,
      "page": 1,
      "limit": 10,
      "total_pages": 5
    },
    "filters": {
      "status": "verified",
      "student_id": "",
      "sort": "created_at",
      "order": "desc"
    }
  }
}
```

#### FR-003: Submit Prestasi
```bash
POST /api/v1/achievements
Authorization: Bearer <token>
Permission: write_achievements

{
  "achievementType": "competition",
  "title": "Juara 1 Hackathon",
  "description": "Deskripsi prestasi",
  "details": { ... },
  "attachments": [],
  "tags": ["hackathon", "programming"]
}
```

#### FR-004: Submit untuk Verifikasi
```bash
POST /api/v1/achievements/:id/submit
Authorization: Bearer <token>
Permission: write_achievements
```

#### FR-005: Hapus Prestasi
```bash
DELETE /api/v1/achievements/:id
Authorization: Bearer <token>
Permission: write_achievements
```

#### FR-006: View Prestasi Mahasiswa Bimbingan (Dosen Wali)
```bash
GET /api/v1/achievements/advisee?page=1&limit=10
Authorization: Bearer <token>
Permission: verify_achievements
```

Response:
```json
{
  "message": "Berhasil mengambil data prestasi mahasiswa bimbingan",
  "data": {
    "achievements": [
      {
        "reference_id": "uuid",
        "achievement_id": "mongo_id",
        "student_id": "NIM123",
        "program_study": "Teknik Informatika",
        "status": "submitted",
        "submitted_at": "2024-12-04T10:00:00Z",
        "achievement": { ... }
      }
    ],
    "pagination": {
      "total": 25,
      "page": 1,
      "limit": 10,
      "total_pages": 3
    }
  }
}
```

#### FR-007: Verify Prestasi (Dosen Wali)
```bash
POST /api/v1/achievements/:id/verify
Authorization: Bearer <token>
Permission: verify_achievements
```

Response:
```json
{
  "message": "Achievement berhasil diverifikasi",
  "data": {
    "achievement_id": "mongo_id",
    "reference_id": "uuid",
    "status": "verified",
    "verified_at": "2024-12-04T10:30:00Z",
    "verified_by": "lecturer_uuid",
    "student_id": "NIM123"
  }
}
```

#### FR-008: Reject Prestasi (Dosen Wali)
```bash
POST /api/v1/achievements/:id/reject
Authorization: Bearer <token>
Permission: verify_achievements
Content-Type: application/json

{
  "rejection_note": "Dokumen pendukung tidak lengkap"
}
```

Response:
```json
{
  "message": "Achievement berhasil ditolak",
  "data": {
    "achievement_id": "mongo_id",
    "reference_id": "uuid",
    "status": "rejected",
    "rejection_note": "Dokumen pendukung tidak lengkap",
    "student_id": "NIM123"
  }
}
```

#### FR-011: Achievement Statistics (Mahasiswa - Own)
```bash
GET /api/v1/achievements/stats/my
Authorization: Bearer <token>
Permission: write_achievements
```

Response:
```json
{
  "message": "Berhasil mengambil statistik prestasi",
  "data": {
    "total": 15,
    "by_type": {
      "competition": 8,
      "academic": 4,
      "organization": 3
    },
    "by_period": {
      "2024-12": 5,
      "2024-11": 7,
      "2024-10": 3
    },
    "by_status": {
      "verified": 10,
      "submitted": 3,
      "draft": 2
    },
    "competition_level_distribution": {
      "international": 2,
      "national": 4,
      "regional": 2
    }
  }
}
```

#### FR-011: Achievement Statistics (Dosen Wali - Advisee)
```bash
GET /api/v1/achievements/stats/advisee
Authorization: Bearer <token>
Permission: verify_achievements
```

Response:
```json
{
  "message": "Berhasil mengambil statistik prestasi mahasiswa bimbingan",
  "data": {
    "total": 45,
    "by_type": {
      "competition": 20,
      "academic": 15,
      "organization": 10
    },
    "by_period": {
      "2024-12": 15,
      "2024-11": 20,
      "2024-10": 10
    },
    "by_status": {
      "verified": 30,
      "submitted": 10,
      "draft": 5
    },
    "competition_level_distribution": {
      "international": 5,
      "national": 10,
      "regional": 5
    },
    "top_students": [
      {
        "student_id": "NIM001",
        "program_study": "Teknik Informatika",
        "count": 8
      },
      {
        "student_id": "NIM002",
        "program_study": "Sistem Informasi",
        "count": 6
      }
    ]
  }
}
```

#### FR-011: Achievement Statistics (Admin - All)
```bash
GET /api/v1/achievements/stats/all
Authorization: Bearer <token>
Permission: read_achievements
```

Response: (Same structure as advisee stats)

### User Management Endpoints (Admin)

#### FR-009: Create User
```bash
POST /api/v1/users
Authorization: Bearer <token>
Permission: manage_users
Content-Type: application/json

{
  "username": "newuser",
  "full_name": "New User",
  "email": "newuser@example.com",
  "password": "password123",
  "role_id": "role-uuid"
}
```

#### FR-009: List Users
```bash
GET /api/v1/users?page=1&limit=10
Authorization: Bearer <token>
Permission: manage_users
```

#### FR-009: Get User Detail
```bash
GET /api/v1/users/:id
Authorization: Bearer <token>
Permission: manage_users
```

#### FR-009: Update User
```bash
PUT /api/v1/users/:id
Authorization: Bearer <token>
Permission: manage_users
Content-Type: application/json

{
  "username": "updateduser",
  "full_name": "Updated Name",
  "email": "updated@example.com",
  "password": "newpassword123"
}
```

#### FR-009: Delete User
```bash
DELETE /api/v1/users/:id
Authorization: Bearer <token>
Permission: manage_users
```

#### FR-009: Assign Role
```bash
PUT /api/v1/users/:id/role
Authorization: Bearer <token>
Permission: manage_users
Content-Type: application/json

{
  "role_id": "role-uuid"
}
```

#### FR-009: Set Student Profile
```bash
POST /api/v1/users/:id/student
Authorization: Bearer <token>
Permission: manage_users
Content-Type: application/json

{
  "student_id": "NIM123",
  "program_study": "Teknik Informatika",
  "academic_year": "2023/2024",
  "advisor_id": "lecturer-uuid"
}
```

#### FR-009: Set Lecturer Profile
```bash
POST /api/v1/users/:id/lecturer
Authorization: Bearer <token>
Permission: manage_users
Content-Type: application/json

{
  "lecturer_id": "NIDN123",
  "department": "Fakultas Teknik"
}
```

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

### Unit Testing

Project ini dilengkapi dengan unit tests untuk memastikan kualitas kode.

#### Run All Tests
```bash
go test ./... -v
```

#### Run Specific Package Tests
```bash
# Test middleware
go test ./Domain/middleware/... -v

# Test services (validation only)
go test ./Domain/service -run "TestSubmitAchievementService|TestRejectAchievementService" -v

# Test user service
go test ./Domain/service -run "TestCreateUserService|TestAssignRoleService" -v
```

#### Test Coverage
```bash
go test ./... -cover
```

#### Test Files
- `Domain/service/AuthService_test.go` - Tests untuk authentication service
- `Domain/service/AchievementService_test.go` - Tests untuk achievement service
- `Domain/service/UserService_test.go` - Tests untuk user management service
- `Domain/middleware/TokenMiddleware_test.go` - Tests untuk JWT middleware

#### Testing Strategy
- **Unit Tests**: Test individual functions dengan mock dependencies
- **Validation Tests**: Test input validation tanpa database
- **Middleware Tests**: Test authentication dan authorization logic
- **Integration Tests**: (Future) Test dengan real database connections

#### Test Examples

**Middleware Test - JWT Authentication:**
```go
func TestJWTAuth_MissingToken(t *testing.T) {
    app := fiber.New()
    app.Use(JWTAuth())
    app.Get("/protected", func(c *fiber.Ctx) error {
        return c.SendString("Protected route")
    })
    
    req := httptest.NewRequest("GET", "/protected", nil)
    resp, err := app.Test(req)
    
    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}
```

**Service Test - Input Validation:**
```go
func TestSubmitAchievementService_InvalidAchievementType(t *testing.T) {
    app := fiber.New()
    app.Post("/achievements", SubmitAchievementService)
    
    achievementData := map[string]interface{}{
        "achievementType": "invalid_type",
        "title":           "Test Achievement",
    }
    body, _ := json.Marshal(achievementData)
    
    req := httptest.NewRequest("POST", "/achievements", bytes.NewBuffer(body))
    resp, err := app.Test(req)
    
    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
```

## 📚 Documentation Files

- `DATABASE_SCHEMA.md` - Skema database lengkap
- `TESTING.md` - Dokumentasi unit testing lengkap
- `PERBAIKAN.md` - Detail perbaikan yang dilakukan
- `RINGKASAN_PERBAIKAN.md` - Ringkasan perbaikan

## 🎯 Features

- ✅ Authentication & Authorization (JWT)
- ✅ Role-based Access Control (RBAC)
- ✅ Token Blacklist (Logout)
- ✅ User Management
- ✅ Achievement Management (PostgreSQL + MongoDB)
- ✅ Flexible Achievement Schema (MongoDB)
- ✅ FR-003: Submit Prestasi (Mahasiswa)
- ✅ FR-004: Submit untuk Verifikasi (Mahasiswa)
- ✅ FR-005: Hapus Prestasi (Mahasiswa)
- ✅ FR-006: View Prestasi Mahasiswa Bimbingan (Dosen Wali)
- ✅ FR-007: Verify Prestasi (Dosen Wali)
- ✅ FR-008: Reject Prestasi (Dosen Wali)
- ✅ FR-009: Manage Users - CRUD, Assign Role, Set Profile (Admin)
- ✅ FR-010: View All Achievements - Filters, Sorting, Pagination (Admin)
- ✅ FR-011: Achievement Statistics - By Type, Period, Top Students (All Roles)

## 📄 License

UAS Project - Backend Lanjutan