package repository

import (
	"GOLANG/Domain/config"
	model "GOLANG/Domain/model/Postgresql"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

func GetUserByEmail(email string) (*model.Users, error) {
	var user model.Users
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users 
		WHERE email = $1
	`
	err := config.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*model.Users, error) {
	var user model.Users
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users 
		WHERE username = $1
	`
	err := config.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByID(userID uuid.UUID) (*model.Users, error) {
	var user model.Users
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users 
		WHERE id = $1
	`
	err := config.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}
	return &user, nil
}
