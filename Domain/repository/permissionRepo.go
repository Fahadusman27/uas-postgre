package repository

import (
	"GOLANG/Domain/config"

	"github.com/google/uuid"
)

func GetPermissionsByRoleID(roleID uuid.UUID) ([]string, error) {
    var permissions []string

    query := `
        SELECT p.name
        FROM permissions p
        JOIN role_permissions rp ON p.id = rp.permission_id
        WHERE rp.role_id = $1
    `

    rows, err := config.DB.Query(query, roleID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            return nil, err
        }
        permissions = append(permissions, name)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return permissions, nil
}