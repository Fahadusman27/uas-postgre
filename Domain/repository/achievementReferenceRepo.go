package repository

import (
	"GOLANG/Domain/config"
	model "GOLANG/Domain/model/Postgresql"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// CreateAchievementReference menyimpan reference ke PostgreSQL
func CreateAchievementReference(ref *model.AchievementReferences) error {
	query := `
		INSERT INTO achievement_references 
		(id, student_id, mongo_achievement_id, status, submitted_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	ref.ID = uuid.New()
	ref.CreatedAt = now
	ref.UpdatedAt = now

	_, err := config.DB.Exec(
		query,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
		ref.SubmittedAt,
		ref.CreatedAt,
		ref.UpdatedAt,
	)

	return err
}

// GetAchievementReferenceByID mengambil reference berdasarkan ID
func GetAchievementReferenceByID(id uuid.UUID) (*model.AchievementReferences, error) {
	var ref model.AchievementReferences
	query := `
		SELECT id, student_id, mongo_achievement_id, status, 
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE id = $1
	`

	err := config.DB.QueryRow(query, id).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("achievement reference tidak ditemukan")
		}
		return nil, err
	}

	return &ref, nil
}

// GetAchievementReferenceByMongoID mengambil reference berdasarkan mongo_achievement_id
func GetAchievementReferenceByMongoID(mongoID string) (*model.AchievementReferences, error) {
	var ref model.AchievementReferences
	query := `
		SELECT id, student_id, mongo_achievement_id, status, 
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE mongo_achievement_id = $1
	`

	err := config.DB.QueryRow(query, mongoID).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("achievement reference tidak ditemukan")
		}
		return nil, err
	}

	return &ref, nil
}

// UpdateAchievementReferenceStatus update status reference
func UpdateAchievementReferenceStatus(id uuid.UUID, status string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := config.DB.Exec(query, status, time.Now(), id)
	return err
}

// UpdateAchievementReference update full reference data
func UpdateAchievementReference(ref *model.AchievementReferences) error {
	query := `
		UPDATE achievement_references
		SET status = $1, submitted_at = $2, verified_at = $3, 
		    verified_by = $4, rejection_note = $5, updated_at = $6
		WHERE id = $7
	`

	ref.UpdatedAt = time.Now()

	_, err := config.DB.Exec(
		query,
		ref.Status,
		ref.SubmittedAt,
		ref.VerifiedAt,
		ref.VerifiedBy,
		ref.RejectionNote,
		ref.UpdatedAt,
		ref.ID,
	)

	return err
}

// DeleteAchievementReference menghapus reference dari PostgreSQL
func DeleteAchievementReference(id uuid.UUID) error {
	query := `DELETE FROM achievement_references WHERE id = $1`
	_, err := config.DB.Exec(query, id)
	return err
}

// GetAchievementReferencesByStudentIDs mengambil references berdasarkan list student IDs dengan pagination
func GetAchievementReferencesByStudentIDs(studentIDs []uuid.UUID, limit, offset int) ([]model.AchievementReferences, int, error) {
	var references []model.AchievementReferences

	// Build query with ANY clause untuk array
	query := `
		SELECT id, student_id, mongo_achievement_id, status, 
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE student_id = ANY($1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	// Convert []uuid.UUID to []string for pq.Array
	studentIDStrings := make([]string, len(studentIDs))
	for i, id := range studentIDs {
		studentIDStrings[i] = id.String()
	}

	rows, err := config.DB.Query(query, "{"+strings.Join(studentIDStrings, ",")+"}", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var ref model.AchievementReferences
		err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		references = append(references, ref)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM achievement_references
		WHERE student_id = ANY($1)
	`
	err = config.DB.QueryRow(countQuery, "{"+strings.Join(studentIDStrings, ",")+"}").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return references, total, nil
}
