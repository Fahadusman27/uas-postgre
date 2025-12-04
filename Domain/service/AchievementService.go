package service

import (
	model "GOLANG/Domain/model/Postgresql"
	mongodb "GOLANG/Domain/model/mongoDB"
	"GOLANG/Domain/repository"
	"log"
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

	// Flow 3: Create notification untuk dosen wali
	// Simple implementation: Log saja untuk sekarang
	log.Printf("Notification: Student %s submitted achievement %s for verification",
		student.StudentID, achievementID)

	// TODO: Implement proper notification system
	// - Save to database
	// - Send email to advisor
	// - Real-time notification

	// Flow 4: Return updated status
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
