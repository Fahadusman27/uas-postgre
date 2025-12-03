package repository

import (
	"GOLANG/Domain/config"
	model "GOLANG/Domain/model/Postgresql"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

// GetStudentByUserID mengambil data student berdasarkan user_id
func GetStudentByUserID(userID uuid.UUID) (*model.Students, error) {
	var student model.Students
	query := `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		WHERE user_id = $1
	`

	err := config.DB.QueryRow(query, userID).Scan(
		&student.ID,
		&student.UserID,
		&student.StudentID,
		&student.ProgramStudy,
		&student.AcademicYear,
		&student.AdvisorID,
		&student.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("student tidak ditemukan")
		}
		return nil, err
	}

	return &student, nil
}

// GetStudentByID mengambil data student berdasarkan student id
func GetStudentByID(studentID uuid.UUID) (*model.Students, error) {
	var student model.Students
	query := `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		WHERE id = $1
	`

	err := config.DB.QueryRow(query, studentID).Scan(
		&student.ID,
		&student.UserID,
		&student.StudentID,
		&student.ProgramStudy,
		&student.AcademicYear,
		&student.AdvisorID,
		&student.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("student tidak ditemukan")
		}
		return nil, err
	}

	return &student, nil
}
