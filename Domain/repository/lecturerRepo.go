package repository

import (
	"GOLANG/Domain/config"
	model "GOLANG/Domain/model/Postgresql"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

// GetLecturerByUserID mengambil data lecturer berdasarkan user_id
func GetLecturerByUserID(userID uuid.UUID) (*model.Lecturers, error) {
	var lecturer model.Lecturers
	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE user_id = $1
	`

	err := config.DB.QueryRow(query, userID).Scan(
		&lecturer.ID,
		&lecturer.UserID,
		&lecturer.LecturerID,
		&lecturer.Department,
		&lecturer.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("lecturer tidak ditemukan")
		}
		return nil, err
	}

	return &lecturer, nil
}

// GetLecturerByID mengambil data lecturer berdasarkan lecturer id
func GetLecturerByID(lecturerID uuid.UUID) (*model.Lecturers, error) {
	var lecturer model.Lecturers
	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE id = $1
	`

	err := config.DB.QueryRow(query, lecturerID).Scan(
		&lecturer.ID,
		&lecturer.UserID,
		&lecturer.LecturerID,
		&lecturer.Department,
		&lecturer.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("lecturer tidak ditemukan")
		}
		return nil, err
	}

	return &lecturer, nil
}
