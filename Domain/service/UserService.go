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
// @Summary Create new user
// @Description Create new user with role assignment (Admin)
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} map[string]interface{} "User created"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 409 {object} map[string]interface{} "Conflict - username/email exists"
// @Router /api/v1/users [post]
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
// @Summary List all users
// @Description Get list of all users with pagination (Admin)
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Router /api/v1/users [get]
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
// @Summary Get user detail
// @Description Get detailed information of a user including profiles (Admin)
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User UUID"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Router /api/v1/users/{id} [get]
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
// @Summary Update user
// @Description Update user information (Admin)
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User UUID"
// @Param user body UpdateUserRequest true "User data"
// @Success 200 {object} map[string]interface{} "Updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Router /api/v1/users/{id} [put]
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User berhasil diupdate",
		"data":    user,
	})
}

func UpdateUserPasswordService(c *fiber.Ctx) error {
    userID := c.Params("id")
    userUUID, err := uuid.Parse(userID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid user ID",
        })
    }

    var req model.Login
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    if req.Password == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Password tidak boleh kosong",
        })
    }

    _, err = repository.GetUserByIDWithDetails(userUUID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "User tidak ditemukan",
        })
    }

    err = repository.UpdateUserPassword(userUUID, req.Password)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Gagal update password",
            "errors": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Password berhasil diupdate",
    })
}

// DeleteUserService - FR-009: Delete User
// @Summary Delete user
// @Description Delete user and associated profiles (Admin)
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User UUID"
// @Success 200 {object} map[string]interface{} "Deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /api/v1/users/{id} [delete]
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

// SetStudentProfileRequest DTO untuk set student profile
type SetStudentProfileRequest struct {
	StudentID    string `json:"student_id"`
	ProgramStudy string `json:"program_study"`
	AcademicYear string `json:"academic_year"`
	AdvisorID    string `json:"advisor_id,omitempty"`
}

// SetStudentProfileService - FR-009: Set Student Profile
// @Summary Set student profile
// @Description Create or update student profile with advisor (Admin)
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User UUID"
// @Param profile body SetStudentProfileRequest true "Student profile data"
// @Success 200 {object} map[string]interface{} "Profile saved"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /api/v1/users/{id}/student [post]
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
// @Summary Set lecturer profile
// @Description Create or update lecturer profile (Admin)
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User UUID"
// @Param profile body SetLecturerProfileRequest true "Lecturer profile data"
// @Success 200 {object} map[string]interface{} "Profile saved"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /api/v1/users/{id}/lecturer [post]
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
