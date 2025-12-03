# FR-004: Submit untuk Verifikasi - Implementation Plan

## Requirement
**Deskripsi**: Mahasiswa submit prestasi draft untuk diverifikasi  
**Actor**: Mahasiswa  
**Precondition**: Prestasi berstatus 'draft'

## Flow
1. Mahasiswa submit prestasi
2. Update status menjadi 'submitted'
3. Create notification untuk dosen wali
4. Return updated status

---

## File yang Sudah Ada (Tidak Perlu Dibuat Baru)

### ✅ Model
1. **`Domain/model/mongoDB/Achievements.go`** - Model Achievement
2. **`Domain/model/Postgresql/achievement_references.go`** - Model AchievementReferences
3. **`Domain/model/Postgresql/Stundent.go`** - Model Students (untuk ambil advisor_id)

### ✅ Repository
1. **`Domain/repository/achievementRepo.go`** - CRUD Achievement (MongoDB)
2. **`Domain/repository/achievementReferenceRepo.go`** - CRUD Reference (PostgreSQL)
3. **`Domain/repository/studentRepo.go`** - Query student data

### ✅ Middleware
1. **`Domain/middleware/TokenMiddleware.go`** - JWT authentication
2. **`Domain/middleware/PermissionMiddleware.go`** - RBAC middleware

### ✅ Service
1. **`Domain/service/AchievementService.go`** - Sudah ada SubmitAchievementService

---

## File yang Perlu DIBUAT/DIMODIFIKASI

### 1. 🆕 BUAT BARU: `Domain/model/Postgresql/Notifications.go` (OPTIONAL)
**Fungsi**: Model untuk notifications (jika ingin simpan ke database)

**Catatan**: Untuk saat ini, kita bisa skip database notification dan langsung implement simple notification system atau log saja. Untuk production bisa pakai:
- Database table `notifications`
- Real-time notification (WebSocket/SSE)
- Email notification
- Push notification

**Rekomendasi**: Skip untuk sekarang, fokus ke flow utama dulu.

---

### 2. 🔄 MODIFIKASI: `Domain/service/AchievementService.go`
**Fungsi**: Tambahkan function `SubmitForVerificationService()`

**Isi**:
```go
// SubmitForVerificationService - FR-004: Submit untuk Verifikasi
func SubmitForVerificationService(c *fiber.Ctx) error {
    // 1. Get achievement_id dari URL parameter
    achievementID := c.Params("id")
    
    // 2. Get user_id dari JWT context
    userID := c.Locals("id").(string)
    
    // 3. Validasi: Cek apakah achievement ada
    // 4. Validasi: Cek apakah user adalah pemilik achievement
    // 5. Validasi: Cek apakah status masih 'draft'
    // 6. Update status menjadi 'submitted'
    // 7. Set submitted_at timestamp
    // 8. Create notification untuk dosen wali (optional)
    // 9. Return updated status
}
```

---

### 3. 🔄 MODIFIKASI: `Domain/route/AchievementRoute.go`
**Fungsi**: Uncomment route untuk submit

**Perubahan**:
```go
// POST /api/v1/achievements/:id/submit - Submit for verification
achievements.Post("/:id/submit",
    middleware.RequirePermission("write_achievements"),
    service.SubmitForVerificationService)
```

---

### 4. 🆕 BUAT BARU (OPTIONAL): `Domain/repository/notificationRepo.go`
**Fungsi**: Repository untuk create notification

**Isi**:
```go
package repository

// CreateNotification membuat notification untuk dosen wali
func CreateNotification(lecturerID uuid.UUID, message string, achievementID string) error {
    // Simple implementation: Log saja untuk sekarang
    log.Printf("Notification for lecturer %s: %s (Achievement: %s)", 
        lecturerID, message, achievementID)
    return nil
}
```

**Catatan**: Untuk sekarang bisa simple log saja. Nanti bisa upgrade ke:
- Database notification
- Email notification
- Real-time notification

---

## Summary: File yang Perlu Dibuat/Dimodifikasi

### File yang Perlu DIMODIFIKASI (2 files):
1. ✅ **`Domain/service/AchievementService.go`** - Tambah `SubmitForVerificationService()`
2. ✅ **`Domain/route/AchievementRoute.go`** - Uncomment route submit

### File OPTIONAL (2 files):
3. ⏳ **`Domain/model/Postgresql/Notifications.go`** - Model notification (optional)
4. ⏳ **`Domain/repository/notificationRepo.go`** - Notification repository (optional)

---

## Implementation Strategy

### Phase 1: Core Functionality (MINIMAL)
**File yang perlu dimodifikasi**: 2 files
1. Tambah function di `AchievementService.go`
2. Uncomment route di `AchievementRoute.go`

**Notification**: Simple log saja

### Phase 2: Notification System (OPTIONAL)
**File yang perlu dibuat**: 2 files
1. Model `Notifications.go`
2. Repository `notificationRepo.go`

**Notification**: Simpan ke database

### Phase 3: Advanced (FUTURE)
- Email notification
- Real-time notification (WebSocket)
- Push notification

---

## Detailed Implementation

### 1. Tambah Function di `AchievementService.go`

```go
// SubmitForVerificationService - FR-004: Submit untuk Verifikasi
func SubmitForVerificationService(c *fiber.Ctx) error {
    // Flow 1: Mahasiswa submit prestasi
    achievementID := c.Params("id")
    
    // Validasi achievement ID
    objectID, err := primitive.ObjectIDFromHex(achievementID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid achievement ID",
        })
    }
    
    // Get user_id dari JWT context
    userID := c.Locals("id").(string)
    userUUID, err := uuid.Parse(userID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid user ID",
        })
    }
    
    // Cek apakah user adalah mahasiswa
    student, err := repository.GetStudentByUserID(userUUID)
    if err != nil {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "User bukan mahasiswa",
        })
    }
    
    // Ambil achievement reference dari PostgreSQL
    reference, err := repository.GetAchievementReferenceByMongoID(achievementID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Achievement tidak ditemukan",
        })
    }
    
    // Validasi: Cek apakah user adalah pemilik achievement
    if reference.StudentID != student.ID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Anda tidak memiliki akses ke achievement ini",
        })
    }
    
    // Precondition: Cek apakah status masih 'draft'
    if reference.Status != "draft" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Achievement hanya bisa di-submit jika berstatus draft",
            "current_status": reference.Status,
        })
    }
    
    // Flow 2: Update status menjadi 'submitted'
    now := time.Now()
    reference.Status = "submitted"
    reference.SubmittedAt = &now
    
    err = repository.UpdateAchievementReference(reference)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Gagal update status achievement",
        })
    }
    
    // Flow 3: Create notification untuk dosen wali
    // Simple implementation: Log saja
    log.Printf("Notification: Student %s submitted achievement %s for verification", 
        student.StudentID, achievementID)
    
    // TODO: Implement proper notification system
    // - Save to database
    // - Send email
    // - Real-time notification
    
    // Flow 4: Return updated status
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Achievement berhasil di-submit untuk verifikasi",
        "data": fiber.Map{
            "achievement_id": achievementID,
            "reference_id": reference.ID,
            "status": reference.Status,
            "submitted_at": reference.SubmittedAt,
        },
    })
}
```

---

### 2. Uncomment Route di `AchievementRoute.go`

```go
// POST /api/v1/achievements/:id/submit - Submit for verification
achievements.Post("/:id/submit",
    middleware.RequirePermission("write_achievements"),
    service.SubmitForVerificationService)
```

---

## Testing Flow

### 1. Create Achievement (Draft)
```bash
TOKEN="mahasiswa-token"
RESPONSE=$(curl -s -X POST http://localhost:4000/api/v1/achievements \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "achievementType": "competition",
    "title": "Test Achievement",
    "description": "Test"
  }')

ACHIEVEMENT_ID=$(echo "$RESPONSE" | jq -r '.data.achievement_id')
echo "Achievement ID: $ACHIEVEMENT_ID"
```

### 2. Submit for Verification
```bash
curl -X POST http://localhost:4000/api/v1/achievements/$ACHIEVEMENT_ID/submit \
  -H "Authorization: Bearer $TOKEN"
```

**Expected Response**:
```json
{
  "message": "Achievement berhasil di-submit untuk verifikasi",
  "data": {
    "achievement_id": "507f1f77bcf86cd799439011",
    "reference_id": "uuid",
    "status": "submitted",
    "submitted_at": "2024-12-04T10:00:00Z"
  }
}
```

### 3. Try Submit Again (Should Fail)
```bash
curl -X POST http://localhost:4000/api/v1/achievements/$ACHIEVEMENT_ID/submit \
  -H "Authorization: Bearer $TOKEN"
```

**Expected Response**:
```json
{
  "error": "Achievement hanya bisa di-submit jika berstatus draft",
  "current_status": "submitted"
}
```

---

## Validation Rules

1. ✅ Achievement ID harus valid
2. ✅ User harus mahasiswa
3. ✅ User harus pemilik achievement
4. ✅ Status harus 'draft' (precondition)
5. ✅ Update status ke 'submitted'
6. ✅ Set submitted_at timestamp
7. ⏳ Create notification (simple log untuk sekarang)

---

## Status Flow

```
draft → submitted → verified
                 ↘ rejected
```

**FR-004 handles**: `draft` → `submitted`

---

## Error Responses

### 400 Bad Request - Invalid ID
```json
{
  "error": "Invalid achievement ID"
}
```

### 400 Bad Request - Wrong Status
```json
{
  "error": "Achievement hanya bisa di-submit jika berstatus draft",
  "current_status": "submitted"
}
```

### 403 Forbidden - Not Owner
```json
{
  "error": "Anda tidak memiliki akses ke achievement ini"
}
```

### 404 Not Found
```json
{
  "error": "Achievement tidak ditemukan"
}
```

---

## Next Steps After FR-004

1. FR-005: Verify Achievement (Dosen)
2. FR-006: Reject Achievement (Dosen)
3. Implement proper notification system
4. Add email notification
5. Add real-time notification

---

## Kesimpulan

**Untuk FR-004, hanya perlu modifikasi 2 files:**
1. ✅ `Domain/service/AchievementService.go` - Tambah function
2. ✅ `Domain/route/AchievementRoute.go` - Uncomment route

**Notification**: Simple log saja untuk sekarang (bisa upgrade nanti)

**Total effort**: Minimal, karena semua infrastructure sudah ada!
