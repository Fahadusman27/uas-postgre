package service

import (
	model "GOLANG/Domain/model/Postgresql"
	mongodb "GOLANG/Domain/model/mongoDB"
	"GOLANG/Domain/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SubmitAchievementRequest DTO untuk request submit prestasi
type SubmitAchievementRequest struct {
	AchievementType string                     `json:"achievementType"`
	Title           string                     `json:"title"`
	Description     string                     `json:"description"`
	Details         mongodb.AchievementDetails `json:"details"`
	CustomFields    map[string]any             `json:"customFields,omitempty"`
	Attachments     []mongodb.Attachment       `json:"attachments"`
	Tags            []string                   `json:"tags"`
}

// SubmitAchievementService - Flow submit prestasi (FR-003)
func SubmitAchievementService(c *fiber.Ctx) error {
	// Flow 1: Mahasiswa mengisi data prestasi
	var req SubmitAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validasi input wajib
	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title wajib diisi",
		})
	}

	if req.AchievementType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Achievement type wajib diisi",
		})
	}

	// Validasi achievement type
	validTypes := map[string]bool{
		"academic":      true,
		"competition":   true,
		"organization":  true,
		"publication":   true,
		"certification": true,
		"other":         true,
	}

	if !validTypes[req.AchievementType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Achievement type tidak valid. Pilihan: academic, competition, organization, publication, certification, other",
		})
	}

	// Ambil user_id dari context (dari JWT middleware)
	userID := c.Locals("id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Precondition: Cek apakah user adalah mahasiswa
	student, err := repository.GetStudentByUserID(userUUID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User bukan mahasiswa atau data mahasiswa tidak ditemukan",
		})
	}

	// Flow 2: Mahasiswa upload dokumen pendukung (sudah ada di req.Attachments)
	// Validasi attachments jika ada
	if req.Attachments == nil {
		req.Attachments = []mongodb.Attachment{}
	}

	// Flow 3: Sistem simpan ke MongoDB (achievement) dan PostgreSQL (reference)

	// 3a. Buat achievement object untuk MongoDB
	achievement := &mongodb.Achievement{
		StudentID:       student.ID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		CustomFields:    req.CustomFields,
		Attachments:     req.Attachments,
		Tags:            req.Tags,
		Points:          0, // Default 0, bisa dihitung nanti berdasarkan rules
	}

	// 3b. Simpan ke MongoDB
	savedAchievement, err := repository.CreateAchievement(achievement)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Gagal menyimpan achievement ke database",
			"details": err.Error(),
		})
	}

	// Flow 4: Status awal: 'draft'
	// 3c. Buat reference di PostgreSQL dengan status 'draft'
	reference := &model.AchievementReferences{
		StudentID:          student.ID,
		MongoAchievementID: savedAchievement.ID.Hex(),
		Status:             "draft", // Status awal: draft
		SubmittedAt:        nil,     // Belum di-submit untuk verifikasi
	}

	err = repository.CreateAchievementReference(reference)
	if err != nil {
		// Rollback: Hapus achievement dari MongoDB jika gagal simpan reference
		_ = repository.DeleteAchievement(savedAchievement.ID)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Gagal menyimpan reference ke database",
			"details": err.Error(),
		})
	}

	// Flow 5: Return achievement data
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Prestasi berhasil disimpan sebagai draft",
		"data": fiber.Map{
			"achievement_id": savedAchievement.ID.Hex(),
			"reference_id":   reference.ID,
			"status":         reference.Status,
			"student_id":     student.ID,
			"achievement":    savedAchievement,
		},
	})
}

// SubmitForVerificationService - FR-004: Submit untuk Verifikasi
func SubmitForVerificationService(c *fiber.Ctx) error {
	// Flow 1: Mahasiswa submit prestasi
	achievementID := c.Params("id")

	// Validasi achievement ID
	_, err := primitive.ObjectIDFromHex(achievementID)
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
			"error":          "Achievement hanya bisa di-submit jika berstatus draft",
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

	// Flow 3: Return updated status
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Achievement berhasil di-submit untuk verifikasi",
		"data": fiber.Map{
			"achievement_id": achievementID,
			"reference_id":   reference.ID,
			"status":         reference.Status,
			"submitted_at":   reference.SubmittedAt,
		},
	})
}

// DeleteAchievementService - FR-005: Hapus Prestasi
func DeleteAchievementService(c *fiber.Ctx) error {
	// Get achievement_id dari URL parameter
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
			"error":          "Hanya achievement dengan status draft yang bisa dihapus",
			"current_status": reference.Status,
		})
	}

	// Flow 1: Soft delete data di MongoDB
	err = repository.SoftDeleteAchievement(objectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus achievement",
		})
	}

	// Flow 2: Delete reference di PostgreSQL
	err = repository.DeleteAchievementReference(reference.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal update reference",
		})
	}

	// Flow 3: Return success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Achievement berhasil dihapus",
		"data": fiber.Map{
			"achievement_id": achievementID,
			"reference_id":   reference.ID,
		},
	})
}

// GetAdviseeAchievementsService - FR-006: View Prestasi Mahasiswa Bimbingan
func GetAdviseeAchievementsService(c *fiber.Ctx) error {
	// Get user_id dari JWT context
	userID := c.Locals("id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Precondition: Cek apakah user adalah dosen
	lecturer, err := repository.GetLecturerByUserID(userUUID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User bukan dosen atau data dosen tidak ditemukan",
		})
	}

	// Flow 1: Get list student IDs dari tabel students where advisor_id
	students, err := repository.GetStudentsByAdvisorID(lecturer.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data mahasiswa bimbingan",
		})
	}

	// Jika tidak ada mahasiswa bimbingan
	if len(students) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Tidak ada mahasiswa bimbingan",
			"data": fiber.Map{
				"achievements": []interface{}{},
				"pagination": fiber.Map{
					"total":       0,
					"page":        1,
					"limit":       10,
					"total_pages": 0,
				},
			},
		})
	}

	// Extract student IDs
	studentIDs := make([]uuid.UUID, len(students))
	for i, student := range students {
		studentIDs[i] = student.ID
	}

	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Flow 2: Get achievements references dengan filter student_ids
	references, total, err := repository.GetAchievementReferencesByStudentIDs(studentIDs, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data achievement references",
		})
	}

	// Jika tidak ada achievements
	if len(references) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Tidak ada prestasi mahasiswa bimbingan",
			"data": fiber.Map{
				"achievements": []interface{}{},
				"pagination": fiber.Map{
					"total":       total,
					"page":        page,
					"limit":       limit,
					"total_pages": 0,
				},
			},
		})
	}

	// Extract mongo achievement IDs
	mongoIDs := make([]string, len(references))
	for i, ref := range references {
		mongoIDs[i] = ref.MongoAchievementID
	}

	// Flow 3: Fetch detail dari MongoDB
	achievements, err := repository.GetAchievementsByMongoIDs(mongoIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil detail achievements dari MongoDB",
		})
	}

	// Create map untuk lookup achievement by mongo ID
	achievementMap := make(map[string]*mongodb.Achievement)
	for i := range achievements {
		achievementMap[achievements[i].ID.Hex()] = &achievements[i]
	}

	// Create map untuk lookup student by ID
	studentMap := make(map[uuid.UUID]*model.Students)
	for i := range students {
		studentMap[students[i].ID] = &students[i]
	}

	// Flow 4: Combine data dan return list dengan pagination
	type AchievementResponse struct {
		ReferenceID   uuid.UUID            `json:"reference_id"`
		AchievementID string               `json:"achievement_id"`
		StudentID     string               `json:"student_id"`
		StudentName   string               `json:"student_name"`
		ProgramStudy  string               `json:"program_study"`
		Status        string               `json:"status"`
		SubmittedAt   *time.Time           `json:"submitted_at"`
		VerifiedAt    *time.Time           `json:"verified_at"`
		Achievement   *mongodb.Achievement `json:"achievement"`
		CreatedAt     time.Time            `json:"created_at"`
	}

	results := make([]AchievementResponse, 0, len(references))
	for _, ref := range references {
		achievement := achievementMap[ref.MongoAchievementID]
		student := studentMap[ref.StudentID]

		if achievement != nil && student != nil {
			results = append(results, AchievementResponse{
				ReferenceID:   ref.ID,
				AchievementID: ref.MongoAchievementID,
				StudentID:     student.StudentID,
				StudentName:   "", // TODO: Get from users table if needed
				ProgramStudy:  student.ProgramStudy,
				Status:        ref.Status,
				SubmittedAt:   ref.SubmittedAt,
				VerifiedAt:    ref.VerifiedAt,
				Achievement:   achievement,
				CreatedAt:     ref.CreatedAt,
			})
		}
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data prestasi mahasiswa bimbingan",
		"data": fiber.Map{
			"achievements": results,
			"pagination": fiber.Map{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": totalPages,
			},
		},
	})
}

// VerifyAchievementService - FR-007: Verify Prestasi
func VerifyAchievementService(c *fiber.Ctx) error {
	// Get achievement_id dari URL parameter
	achievementID := c.Params("id")

	// Validasi achievement ID
	_, err := primitive.ObjectIDFromHex(achievementID)
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

	// Flow 1: Validasi user adalah dosen wali
	lecturer, err := repository.GetLecturerByUserID(userUUID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User bukan dosen atau data dosen tidak ditemukan",
		})
	}

	// Get achievement reference dari PostgreSQL
	reference, err := repository.GetAchievementReferenceByMongoID(achievementID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Achievement tidak ditemukan",
		})
	}

	// Precondition: Cek apakah status 'submitted'
	if reference.Status != "submitted" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":          "Achievement hanya bisa diverifikasi jika berstatus submitted",
			"current_status": reference.Status,
		})
	}

	// Validasi: Cek apakah mahasiswa adalah anak bimbingan dosen ini
	student, err := repository.GetStudentByID(reference.StudentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Data mahasiswa tidak ditemukan",
		})
	}

	// Cek apakah advisor_id mahasiswa sama dengan lecturer ID
	if student.AdvisorID != lecturer.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Anda tidak memiliki akses untuk memverifikasi prestasi mahasiswa ini",
		})
	}

	// Flow 2: Dosen approve prestasi
	// Flow 3: Update status menjadi 'verified'
	now := time.Now()
	reference.Status = "verified"
	reference.VerifiedAt = &now
	reference.VerifiedBy = &lecturer.ID
	reference.RejectionNote = nil // Clear rejection note jika ada

	// Flow 4: Set verified_by dan verified_at (sudah dilakukan di atas)
	err = repository.UpdateAchievementReference(reference)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal update status achievement",
		})
	}

	// Flow 5: Return updated status
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Achievement berhasil diverifikasi",
		"data": fiber.Map{
			"achievement_id": achievementID,
			"reference_id":   reference.ID,
			"status":         reference.Status,
			"verified_at":    reference.VerifiedAt,
			"verified_by":    reference.VerifiedBy,
			"student_id":     student.StudentID,
		},
	})
}

// RejectAchievementRequest DTO untuk request reject prestasi
type RejectAchievementRequest struct {
	RejectionNote string `json:"rejection_note"`
}

// RejectAchievementService - FR-008: Reject Prestasi
func RejectAchievementService(c *fiber.Ctx) error {
	// Get achievement_id dari URL parameter
	achievementID := c.Params("id")

	// Validasi achievement ID
	_, err := primitive.ObjectIDFromHex(achievementID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid achievement ID",
		})
	}

	// Flow 1: Dosen input rejection note
	var req RejectAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validasi rejection note wajib diisi
	if req.RejectionNote == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Rejection note wajib diisi",
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

	// Validasi user adalah dosen wali
	lecturer, err := repository.GetLecturerByUserID(userUUID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User bukan dosen atau data dosen tidak ditemukan",
		})
	}

	// Get achievement reference dari PostgreSQL
	reference, err := repository.GetAchievementReferenceByMongoID(achievementID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Achievement tidak ditemukan",
		})
	}

	// Precondition: Cek apakah status 'submitted'
	if reference.Status != "submitted" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":          "Achievement hanya bisa ditolak jika berstatus submitted",
			"current_status": reference.Status,
		})
	}

	// Validasi: Cek apakah mahasiswa adalah anak bimbingan dosen ini
	student, err := repository.GetStudentByID(reference.StudentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Data mahasiswa tidak ditemukan",
		})
	}

	// Cek apakah advisor_id mahasiswa sama dengan lecturer ID
	if student.AdvisorID != lecturer.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Anda tidak memiliki akses untuk menolak prestasi mahasiswa ini",
		})
	}

	// Flow 2: Update status menjadi 'rejected'
	// Flow 3: Save rejection_note
	reference.Status = "rejected"
	reference.RejectionNote = &req.RejectionNote
	reference.VerifiedAt = nil // Clear verified_at
	reference.VerifiedBy = nil // Clear verified_by

	err = repository.UpdateAchievementReference(reference)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal update status achievement",
		})
	}

	// Flow 4: Return updated status
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Achievement berhasil ditolak",
		"data": fiber.Map{
			"achievement_id": achievementID,
			"reference_id":   reference.ID,
			"status":         reference.Status,
			"rejection_note": reference.RejectionNote,
			"student_id":     student.StudentID,
		},
	})
}

// GetAllAchievementsService - FR-010: View All Achievements (Admin)
func GetAllAchievementsService(c *fiber.Ctx) error {
	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Parse filter parameters
	status := c.Query("status", "")
	studentID := c.Query("student_id", "")
	sortBy := c.Query("sort", "created_at")
	order := c.Query("order", "desc")

	// Flow 1: Get all achievement references dengan filters
	references, total, err := repository.GetAllAchievementReferencesWithFilters(limit, offset, status, studentID, sortBy, order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data achievement references",
		})
	}

	// Jika tidak ada achievements
	if len(references) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Tidak ada prestasi",
			"data": fiber.Map{
				"achievements": []interface{}{},
				"pagination": fiber.Map{
					"total":       total,
					"page":        page,
					"limit":       limit,
					"total_pages": 0,
				},
			},
		})
	}

	// Extract mongo achievement IDs
	mongoIDs := make([]string, len(references))
	studentIDs := make(map[uuid.UUID]bool)
	for i, ref := range references {
		mongoIDs[i] = ref.MongoAchievementID
		studentIDs[ref.StudentID] = true
	}

	// Flow 2: Fetch details dari MongoDB
	achievements, err := repository.GetAchievementsByMongoIDs(mongoIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil detail achievements dari MongoDB",
		})
	}

	// Create map untuk lookup achievement by mongo ID
	achievementMap := make(map[string]*mongodb.Achievement)
	for i := range achievements {
		achievementMap[achievements[i].ID.Hex()] = &achievements[i]
	}

	// Get student data
	studentIDList := make([]uuid.UUID, 0, len(studentIDs))
	for id := range studentIDs {
		studentIDList = append(studentIDList, id)
	}

	studentMap := make(map[uuid.UUID]*model.Students)
	for _, id := range studentIDList {
		student, err := repository.GetStudentByID(id)
		if err == nil {
			studentMap[id] = student
		}
	}

	// Flow 3: Combine data dan return dengan pagination
	type AchievementResponse struct {
		ReferenceID   uuid.UUID            `json:"reference_id"`
		AchievementID string               `json:"achievement_id"`
		StudentID     string               `json:"student_id"`
		ProgramStudy  string               `json:"program_study"`
		Status        string               `json:"status"`
		SubmittedAt   *time.Time           `json:"submitted_at"`
		VerifiedAt    *time.Time           `json:"verified_at"`
		VerifiedBy    *uuid.UUID           `json:"verified_by"`
		RejectionNote *string              `json:"rejection_note"`
		Achievement   *mongodb.Achievement `json:"achievement"`
		CreatedAt     time.Time            `json:"created_at"`
	}

	results := make([]AchievementResponse, 0, len(references))
	for _, ref := range references {
		achievement := achievementMap[ref.MongoAchievementID]
		student := studentMap[ref.StudentID]

		response := AchievementResponse{
			ReferenceID:   ref.ID,
			AchievementID: ref.MongoAchievementID,
			Status:        ref.Status,
			SubmittedAt:   ref.SubmittedAt,
			VerifiedAt:    ref.VerifiedAt,
			VerifiedBy:    ref.VerifiedBy,
			RejectionNote: ref.RejectionNote,
			Achievement:   achievement,
			CreatedAt:     ref.CreatedAt,
		}

		if student != nil {
			response.StudentID = student.StudentID
			response.ProgramStudy = student.ProgramStudy
		}

		results = append(results, response)
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	// Flow 4: Return dengan pagination
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data achievements",
		"data": fiber.Map{
			"achievements": results,
			"pagination": fiber.Map{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": totalPages,
			},
			"filters": fiber.Map{
				"status":     status,
				"student_id": studentID,
				"sort":       sortBy,
				"order":      order,
			},
		},
	})
}

// GetMyAchievementStatsService - FR-011: Achievement Statistics (Mahasiswa - Own)
func GetMyAchievementStatsService(c *fiber.Ctx) error {
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

	// Get all achievement references untuk student ini
	references, err := repository.GetAllAchievementReferences(student.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data achievements",
		})
	}

	if len(references) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Belum ada prestasi",
			"data": fiber.Map{
				"total":                          0,
				"by_type":                        map[string]int{},
				"by_period":                      map[string]int{},
				"by_status":                      map[string]int{},
				"competition_level_distribution": map[string]int{},
			},
		})
	}

	// Extract mongo IDs
	mongoIDs := make([]string, len(references))
	for i, ref := range references {
		mongoIDs[i] = ref.MongoAchievementID
	}

	// Get statistics
	statsByType, _ := repository.GetAchievementStatsByType(mongoIDs)
	statsByPeriod, _ := repository.GetAchievementStatsByPeriod(mongoIDs)
	competitionLevelDist, _ := repository.GetCompetitionLevelDistribution(mongoIDs)

	// Count by status
	statsByStatus := make(map[string]int)
	for _, ref := range references {
		statsByStatus[ref.Status]++
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil statistik prestasi",
		"data": fiber.Map{
			"total":                          len(references),
			"by_type":                        statsByType,
			"by_period":                      statsByPeriod,
			"by_status":                      statsByStatus,
			"competition_level_distribution": competitionLevelDist,
		},
	})
}

// GetAdviseeAchievementStatsService - FR-011: Achievement Statistics (Dosen Wali - Advisee)
func GetAdviseeAchievementStatsService(c *fiber.Ctx) error {
	// Get user_id dari JWT context
	userID := c.Locals("id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Cek apakah user adalah dosen
	lecturer, err := repository.GetLecturerByUserID(userUUID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User bukan dosen",
		})
	}

	// Get list students
	students, err := repository.GetStudentsByAdvisorID(lecturer.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data mahasiswa bimbingan",
		})
	}

	if len(students) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Tidak ada mahasiswa bimbingan",
			"data": fiber.Map{
				"total":                          0,
				"by_type":                        map[string]int{},
				"by_period":                      map[string]int{},
				"by_status":                      map[string]int{},
				"competition_level_distribution": map[string]int{},
				"top_students":                   []interface{}{},
			},
		})
	}

	// Extract student IDs
	studentIDs := make([]uuid.UUID, len(students))
	for i, student := range students {
		studentIDs[i] = student.ID
	}

	// Get all achievement references
	references, total, err := repository.GetAchievementReferencesByStudentIDs(studentIDs, 10000, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data achievements",
		})
	}

	if len(references) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Belum ada prestasi mahasiswa bimbingan",
			"data": fiber.Map{
				"total":                          0,
				"by_type":                        map[string]int{},
				"by_period":                      map[string]int{},
				"by_status":                      map[string]int{},
				"competition_level_distribution": map[string]int{},
				"top_students":                   []interface{}{},
			},
		})
	}

	// Extract mongo IDs
	mongoIDs := make([]string, len(references))
	for i, ref := range references {
		mongoIDs[i] = ref.MongoAchievementID
	}

	// Get statistics
	statsByType, _ := repository.GetAchievementStatsByType(mongoIDs)
	statsByPeriod, _ := repository.GetAchievementStatsByPeriod(mongoIDs)
	competitionLevelDist, _ := repository.GetCompetitionLevelDistribution(mongoIDs)
	statsByStatus, _ := repository.GetAchievementCountByStatus(studentIDs)

	// Get top students
	topStudents, _ := repository.GetTopStudentsByAchievementCount(studentIDs, 10, "verified")

	// Create student map for lookup
	studentMap := make(map[uuid.UUID]*model.Students)
	for i := range students {
		studentMap[students[i].ID] = &students[i]
	}

	// Format top students
	type TopStudentResponse struct {
		StudentID    string `json:"student_id"`
		ProgramStudy string `json:"program_study"`
		Count        int    `json:"count"`
	}

	topStudentsResponse := make([]TopStudentResponse, 0, len(topStudents))
	for _, ts := range topStudents {
		if student, ok := studentMap[ts.StudentID]; ok {
			topStudentsResponse = append(topStudentsResponse, TopStudentResponse{
				StudentID:    student.StudentID,
				ProgramStudy: student.ProgramStudy,
				Count:        ts.Count,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil statistik prestasi mahasiswa bimbingan",
		"data": fiber.Map{
			"total":                          total,
			"by_type":                        statsByType,
			"by_period":                      statsByPeriod,
			"by_status":                      statsByStatus,
			"competition_level_distribution": competitionLevelDist,
			"top_students":                   topStudentsResponse,
		},
	})
}

// GetAllAchievementStatsService - FR-011: Achievement Statistics (Admin - All)
func GetAllAchievementStatsService(c *fiber.Ctx) error {
	// Get all achievement references (no filter)
	references, total, err := repository.GetAllAchievementReferencesWithFilters(100000, 0, "", "", "created_at", "desc")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data achievements",
		})
	}

	if len(references) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Belum ada prestasi",
			"data": fiber.Map{
				"total":                          0,
				"by_type":                        map[string]int{},
				"by_period":                      map[string]int{},
				"by_status":                      map[string]int{},
				"competition_level_distribution": map[string]int{},
				"top_students":                   []interface{}{},
			},
		})
	}

	// Extract mongo IDs and student IDs
	mongoIDs := make([]string, len(references))
	studentIDMap := make(map[uuid.UUID]bool)
	for i, ref := range references {
		mongoIDs[i] = ref.MongoAchievementID
		studentIDMap[ref.StudentID] = true
	}

	studentIDs := make([]uuid.UUID, 0, len(studentIDMap))
	for id := range studentIDMap {
		studentIDs = append(studentIDs, id)
	}

	// Get statistics
	statsByType, _ := repository.GetAchievementStatsByType(mongoIDs)
	statsByPeriod, _ := repository.GetAchievementStatsByPeriod(mongoIDs)
	competitionLevelDist, _ := repository.GetCompetitionLevelDistribution(mongoIDs)
	statsByStatus, _ := repository.GetAchievementCountByStatus(studentIDs)

	// Get top students
	topStudents, _ := repository.GetTopStudentsByAchievementCount(studentIDs, 10, "verified")

	// Get student data for top students
	type TopStudentResponse struct {
		StudentID    string `json:"student_id"`
		ProgramStudy string `json:"program_study"`
		Count        int    `json:"count"`
	}

	topStudentsResponse := make([]TopStudentResponse, 0, len(topStudents))
	for _, ts := range topStudents {
		student, err := repository.GetStudentByID(ts.StudentID)
		if err == nil {
			topStudentsResponse = append(topStudentsResponse, TopStudentResponse{
				StudentID:    student.StudentID,
				ProgramStudy: student.ProgramStudy,
				Count:        ts.Count,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil statistik prestasi",
		"data": fiber.Map{
			"total":                          total,
			"by_type":                        statsByType,
			"by_period":                      statsByPeriod,
			"by_status":                      statsByStatus,
			"competition_level_distribution": competitionLevelDist,
			"top_students":                   topStudentsResponse,
		},
	})
}
