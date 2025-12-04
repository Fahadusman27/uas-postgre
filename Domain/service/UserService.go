package service

import (
	model "GOLANG/Domain/model/Postgresql"
	"GOLANG/Domain/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateUserRequest DTO untuk create user
type CreateUserRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   string `json:"role_id"`
}

// CreateUserService - FR-009: Create User
func CreateUserService(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validasi input
	if req.Username == "" || req.FullName == "" || req.Email == "" || req.Password == "" || req.RoleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Semua field wajib diisi",
		})
	}

	// Parse role_id
	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	// Cek apakah username sudah ada
	existingUser, _ := repository.GetUserByUsername(req.Username)
	if existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Username sudah digunakan",
		})
	}

	// Cek apakah email sudah ada
	existingUser, _ = repository.GetUserByEmail(req.Email)
	if existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email sudah digunakan",
		})
	}

	// Create user
	user := &model.Users{
		Username: req.Username,
		FullName: req.FullName,
		Email:    req.Email,
		RoleID:   roleID,
	}

	err = repository.CreateUser(user, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuat user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User berhasil dibuat",
		"data": fiber.Map{
			"id":         user.ID,
			"username":   user.Username,
			"full_name":  user.FullName,
			"email":      user.Email,
			"role_id":    user.RoleID,
			"created_at": user.CreatedAt,
		},
	})
}

// GetUsersService - FR-009: List Users
func GetUsersService(c *fiber.Ctx) error {
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

	// Get users
	users, total, err := repository.GetAllUsers(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data users",
		})
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data users",
		"data": fiber.Map{
			"users": users,
			"pagination": fiber.Map{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": totalPages,
			},
		},
	})
}

// GetUserDetailService - FR-009: Get User Detail
func GetUserDetailService(c *fiber.Ctx) error {
	userID := c.Params("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get user
	user, err := repository.GetUserByIDWithDetails(userUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User tidak ditemukan",
		})
	}

	// Get student profile if exists
	student, _ := repository.GetStudentByUserID(userUUID)

	// Get lecturer profile if exists
	lecturer, _ := repository.GetLecturerByUserID(userUUID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil detail user",
		"data": fiber.Map{
			"user":     user,
			"student":  student,
			"lecturer": lecturer,
		},
	})
}

// UpdateUserRequest DTO untuk update user
type UpdateUserRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// UpdateUserService - FR-009: Update User
func UpdateUserService(c *fiber.Ctx) error {
	userID := c.Params("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get existing user
	user, err := repository.GetUserByIDWithDetails(userUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User tidak ditemukan",
		})
	}

	// Update fields
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	// Update user
	err = repository.UpdateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal update user",
		})
	}

	// Update password if provided
	if req.Password != "" {
		err = repository.UpdateUserPassword(userUUID, req.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal update password",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User berhasil diupdate",
		"data":    user,
	})
}

// DeleteUserService - FR-009: Delete User
func DeleteUserService(c *fiber.Ctx) error {
	userID := c.Params("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Delete user
	err = repository.DeleteUser(userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User berhasil dihapus",
	})
}

// AssignRoleRequest DTO untuk assign role
type AssignRoleRequest struct {
	RoleID string `json:"role_id"`
}

// AssignRoleService - FR-009: Assign Role
func AssignRoleService(c *fiber.Ctx) error {
	userID := c.Params("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	// Update role
	err = repository.UpdateUserRole(userUUID, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal assign role",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role berhasil di-assign",
		"data": fiber.Map{
			"user_id": userUUID,
			"role_id": roleID,
		},
	})
}

// SetStudentProfileRequest DTO untuk set student profile
type SetStudentProfileRequest struct {
	StudentID    string `json:"student_id"`
	ProgramStudy string `json:"program_study"`
	AcademicYear string `json:"academic_year"`
	AdvisorID    string `json:"advisor_id,omitempty"`
}

// SetStudentProfileService - FR-009: Set Student Profile
func SetStudentProfileService(c *fiber.Ctx) error {
	userID := c.Params("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req SetStudentProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validasi input
	if req.StudentID == "" || req.ProgramStudy == "" || req.AcademicYear == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID, program study, dan academic year wajib diisi",
		})
	}

	// Parse advisor_id if provided
	var advisorID uuid.UUID
	if req.AdvisorID != "" {
		advisorID, err = uuid.Parse(req.AdvisorID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid advisor ID",
			})
		}
	}

	// Check if profile already exists
	exists, err := repository.CheckStudentProfileExists(userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal cek student profile",
		})
	}

	student := &model.Students{
		UserID:       userUUID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
		AdvisorID:    advisorID,
	}

	if exists {
		// Update existing profile
		err = repository.UpdateStudentProfile(student)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal update student profile",
			})
		}
	} else {
		// Create new profile
		err = repository.CreateStudentProfile(student)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal membuat student profile",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Student profile berhasil disimpan",
		"data":    student,
	})
}

// SetLecturerProfileRequest DTO untuk set lecturer profile
type SetLecturerProfileRequest struct {
	LecturerID string `json:"lecturer_id"`
	Department string `json:"department"`
}

// SetLecturerProfileService - FR-009: Set Lecturer Profile
func SetLecturerProfileService(c *fiber.Ctx) error {
	userID := c.Params("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req SetLecturerProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validasi input
	if req.LecturerID == "" || req.Department == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Lecturer ID dan department wajib diisi",
		})
	}

	// Check if profile already exists
	exists, err := repository.CheckLecturerProfileExists(userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal cek lecturer profile",
		})
	}

	lecturer := &model.Lecturers{
		UserID:     userUUID,
		LecturerID: req.LecturerID,
		Department: req.Department,
	}

	if exists {
		// Update existing profile
		err = repository.UpdateLecturerProfile(lecturer)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal update lecturer profile",
			})
		}
	} else {
		// Create new profile
		err = repository.CreateLecturerProfile(lecturer)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal membuat lecturer profile",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Lecturer profile berhasil disimpan",
		"data":    lecturer,
	})
}
