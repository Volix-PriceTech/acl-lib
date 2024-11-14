package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Volix-PriceTech/acl-lib/models"
)

type MySQLStorage struct {
	db *sql.DB
}

func NewMySQLStorage(db *sql.DB) Storage {
	return &MySQLStorage{db: db}
}

func (m *MySQLStorage) MigrateTables() error {
	queries := []string{
		`
        CREATE TABLE IF NOT EXISTS roles (
            id BIGINT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE
        );
        `,
		`
        CREATE TABLE IF NOT EXISTS permissions (
            id BIGINT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE
        );
        `,
		`
        CREATE TABLE IF NOT EXISTS roles_permissions (
            role_id BIGINT NOT NULL,
            permission_id BIGINT NOT NULL,
            PRIMARY KEY (role_id, permission_id),
            FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
            FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
        );
        `,
	}

	for _, query := range queries {
		if _, err := m.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration query: %w", err)
		}
	}
	return nil
}

func (m *MySQLStorage) CreateRole(role *models.Role) error {
	query := "INSERT INTO roles (id, name) VALUES (?, ?)"
	_, err := m.db.Exec(query, role.ID, role.Name)
	return err
}

func (m *MySQLStorage) GetRoleByID(id string) (*models.Role, error) {
	query := "SELECT id, name FROM roles WHERE id = ?"
	row := m.db.QueryRow(query, id)

	var role models.Role
	if err := row.Scan(&role.ID, &role.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (m *MySQLStorage) DeleteRole(id string) error {
	query := "DELETE FROM roles WHERE id = ?"
	result, err := m.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("role not found")
	}
	return nil
}

func (m *MySQLStorage) CreatePermission(permission *models.Permission) error {
	query := "INSERT INTO permissions (id, name) VALUES (?, ?)"
	_, err := m.db.Exec(query, permission.ID, permission.Name)
	return err
}

func (m *MySQLStorage) GetPermissionByID(id string) (*models.Permission, error) {
	query := "SELECT id, name FROM permissions WHERE id = ?"
	row := m.db.QueryRow(query, id)

	var permission models.Permission
	if err := row.Scan(&permission.ID, &permission.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}
	return &permission, nil
}

func (m *MySQLStorage) DeletePermission(id string) error {
	query := "DELETE FROM permissions WHERE id = ?"
	result, err := m.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("permission not found")
	}
	return nil
}

func (m *MySQLStorage) AssignPermissionToRole(roleID, permissionID string) error {
	query := "INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)"
	_, err := m.db.Exec(query, roleID, permissionID)
	return err
}

func (m *MySQLStorage) RemovePermissionFromRole(roleID, permissionID string) error {
	query := "DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?"
	result, err := m.db.Exec(query, roleID, permissionID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("assignment not found")
	}
	return nil
}

func (m *MySQLStorage) GetPermissionsByRole(roleID string) ([]models.Permission, error) {
	query := `
		SELECT p.id, p.name
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?`

	rows, err := m.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var permission models.Permission
		if err := rows.Scan(&permission.ID, &permission.Name); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return permissions, nil
}
