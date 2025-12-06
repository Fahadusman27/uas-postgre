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

func GetAllAchievementReferences(StudentID uuid.UUID) ([]model.AchievementReferences, error) {
	query := `
        SELECT 
            id, 
            student_id, 
            mongo_achievement_id, 
            status, 
            submitted_at, 
            verified_at, 
            verified_by, 
            rejection_note, 
            created_at, 
            updated_at 
        FROM achievement_references 
        WHERE student_id = $1`

	// 2. Gunakan Query (bukan Exec)
	rows, err := config.DB.Query(query, StudentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var references []model.AchievementReferences

	// 3. Loop setiap baris data
	for rows.Next() {
		var ref model.AchievementReferences

		// 4. Scan data ke variable struct
		// PERHATIAN: Urutan Scan harus SAMA PERSIS dengan urutan kolom di SELECT
		err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,   // Pointer otomatis menangani NULL
			&ref.VerifiedAt,    // Pointer otomatis menangani NULL
			&ref.VerifiedBy,    // Pointer otomatis menangani NULL
			&ref.RejectionNote, // Pointer otomatis menangani NULL
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		references = append(references, ref)
	}

	// Cek error setelah iterasi selesai
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return references, nil
}

// GetAllAchievementReferencesWithFilters mengambil semua achievement references dengan filters dan pagination
func GetAllAchievementReferencesWithFilters(limit, offset int, status, studentID, sortBy, order string) ([]model.AchievementReferences, int, error) {
	var references []model.AchievementReferences

	// Build query dengan filters
	query := `
		SELECT id, student_id, mongo_achievement_id, status, 
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE 1=1
	`

	// Count query
	countQuery := `SELECT COUNT(*) FROM achievement_references WHERE 1=1`

	// Build WHERE clause
	var args []interface{}
	argIndex := 1

	// Filter by status
	if status != "" {
		query += " AND status = $" + string(rune(argIndex+'0'))
		countQuery += " AND status = $1"
		args = append(args, status)
		argIndex++
	}

	// Filter by student_id
	if studentID != "" {
		studentUUID, err := uuid.Parse(studentID)
		if err == nil {
			placeholder := "$" + string(rune(argIndex+'0'))
			query += " AND student_id = " + placeholder
			if len(args) == 0 {
				countQuery += " AND student_id = $1"
			} else {
				countQuery += " AND student_id = $2"
			}
			args = append(args, studentUUID)
			argIndex++
		}
	}

	// Sorting
	validSortFields := map[string]bool{
		"created_at":   true,
		"submitted_at": true,
		"verified_at":  true,
		"updated_at":   true,
	}

	if sortBy == "" || !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	query += " ORDER BY " + sortBy + " " + order

	// Pagination
	query += " LIMIT $" + string(rune(argIndex+'0')) + " OFFSET $" + string(rune(argIndex+'1'))
	args = append(args, limit, offset)

	// Execute query
	rows, err := config.DB.Query(query, args...)
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
	countArgs := args[:len(args)-2] // Remove limit and offset
	err = config.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return references, total, nil
}

// GetTopStudentsByAchievementCount mengambil top mahasiswa berdasarkan jumlah prestasi
func GetTopStudentsByAchievementCount(studentIDs []uuid.UUID, limit int, status string) ([]struct {
	StudentID uuid.UUID
	Count     int
}, error) {
	query := `
		SELECT student_id, COUNT(*) as count
		FROM achievement_references
		WHERE student_id = ANY($1)
	`

	// Add status filter if provided
	args := []interface{}{"{" + strings.Join(func() []string {
		strs := make([]string, len(studentIDs))
		for i, id := range studentIDs {
			strs[i] = id.String()
		}
		return strs
	}(), ",") + "}"}

	if status != "" {
		query += " AND status = $2"
		args = append(args, status)
		query += " GROUP BY student_id ORDER BY count DESC LIMIT $3"
		args = append(args, limit)
	} else {
		query += " GROUP BY student_id ORDER BY count DESC LIMIT $2"
		args = append(args, limit)
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type StudentCount struct {
		StudentID uuid.UUID
		Count     int
	}

	var results []StudentCount
	for rows.Next() {
		var sc StudentCount
		err := rows.Scan(&sc.StudentID, &sc.Count)
		if err != nil {
			return nil, err
		}
		results = append(results, sc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Convert to anonymous struct slice
	finalResults := make([]struct {
		StudentID uuid.UUID
		Count     int
	}, len(results))

	for i, r := range results {
		finalResults[i].StudentID = r.StudentID
		finalResults[i].Count = r.Count
	}

	return finalResults, nil
}

// GetAchievementCountByStatus mengambil jumlah prestasi per status
func GetAchievementCountByStatus(studentIDs []uuid.UUID) (map[string]int, error) {
	query := `
		SELECT status, COUNT(*) as count
		FROM achievement_references
		WHERE student_id = ANY($1)
		GROUP BY status
	`

	studentIDStrings := make([]string, len(studentIDs))
	for i, id := range studentIDs {
		studentIDStrings[i] = id.String()
	}

	rows, err := config.DB.Query(query, "{"+strings.Join(studentIDStrings, ",")+"}")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		err := rows.Scan(&status, &count)
		if err != nil {
			return nil, err
		}
		stats[status] = count
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
