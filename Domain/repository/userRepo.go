package repository

import (
	"GOLANG/Domain/config"
	model "GOLANG/Domain/model/Postgresql"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser membuat user baru
func CreateUser(user *model.Users, password string) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (id, username, full_name, email, password, role_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	user.ID = uuid.New()
	user.CreatedAt = time.Now()

	_, err = config.DB.Exec(
		query,
		user.ID,
		user.Username,
		user.FullName,
		user.Email,
		string(hashedPassword),
		user.RoleID,
		user.CreatedAt,
	)

	return err
}

// GetUserByIDWithDetails mengambil user berdasarkan ID (alias untuk menghindari konflik)
func GetUserByIDWithDetails(userID uuid.UUID) (*model.Users, error) {
	var user model.Users
	query := `
		SELECT id, username, full_name, email, role_id, created_at
		FROM users
		WHERE id = $1
	`

	err := config.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.FullName,
		&user.Email,
		&user.RoleID,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}

	return &user, nil
}

// GetAllUsers mengambil semua users dengan pagination
func GetAllUsers(limit, offset int) ([]model.Users, int, error) {
	var users []model.Users

	query := `
		SELECT id, username, full_name, email, role_id, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := config.DB.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.Users
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.FullName,
			&user.Email,
			&user.RoleID,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM users`
	err = config.DB.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser update user data
func UpdateUser(user *model.Users) error {
	query := `
		UPDATE users
		SET username = $1, full_name = $2, email = $3, role_id = $4
		WHERE id = $5
	`

	result, err := config.DB.Exec(
		query,
		user.Username,
		user.FullName,
		user.Email,
		user.RoleID,
		user.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user tidak ditemukan")
	}

	return nil
}

// UpdateUserPassword update password user
func UpdateUserPassword(userID uuid.UUID, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `UPDATE users SET password_hash = $1 WHERE id = $2`
	_, err = config.DB.Exec(query, string(hashedPassword), userID)
	return err
}

// DeleteUser menghapus user
func DeleteUser(userID uuid.UUID) error {
	// Hard delete - hapus juga profile student/lecturer jika ada
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete student profile if exists
	_, err = tx.Exec("DELETE FROM students WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	// Delete lecturer profile if exists
	_, err = tx.Exec("DELETE FROM lecturers WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	// Delete user
	result, err := tx.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user tidak ditemukan")
	}

	return tx.Commit()
}

// CreateStudentProfile membuat profile student
func CreateStudentProfile(student *model.Students) error {
	query := `
		INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	student.ID = uuid.New()
	student.CreatedAt = time.Now()

	_, err := config.DB.Exec(
		query,
		student.ID,
		student.UserID,
		student.StudentID,
		student.ProgramStudy,
		student.AcademicYear,
		student.AdvisorID,
		student.CreatedAt,
	)

	return err
}

// UpdateStudentProfile update profile student termasuk advisor
func UpdateStudentProfile(student *model.Students) error {
	query := `
		UPDATE students
		SET student_id = $1, program_study = $2, academic_year = $3, advisor_id = $4
		WHERE user_id = $5
	`

	result, err := config.DB.Exec(
		query,
		student.StudentID,
		student.ProgramStudy,
		student.AcademicYear,
		student.AdvisorID,
		student.UserID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("student profile tidak ditemukan")
	}

	return nil
}

// CreateLecturerProfile membuat profile lecturer
func CreateLecturerProfile(lecturer *model.Lecturers) error {
	query := `
		INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	lecturer.ID = uuid.New()
	lecturer.CreatedAt = time.Now()

	_, err := config.DB.Exec(
		query,
		lecturer.ID,
		lecturer.UserID,
		lecturer.LecturerID,
		lecturer.Department,
		lecturer.CreatedAt,
	)

	return err
}

// UpdateLecturerProfile update profile lecturer
func UpdateLecturerProfile(lecturer *model.Lecturers) error {
	query := `
		UPDATE lecturers
		SET lecturer_id = $1, department = $2
		WHERE user_id = $3
	`

	result, err := config.DB.Exec(
		query,
		lecturer.LecturerID,
		lecturer.Department,
		lecturer.UserID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("lecturer profile tidak ditemukan")
	}

	return nil
}

// CheckStudentProfileExists cek apakah student profile sudah ada
func CheckStudentProfileExists(userID uuid.UUID) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM students WHERE user_id = $1`
	err := config.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckLecturerProfileExists cek apakah lecturer profile sudah ada
func CheckLecturerProfileExists(userID uuid.UUID) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM lecturers WHERE user_id = $1`
	err := config.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
