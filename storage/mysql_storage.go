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
        CREATE TABLE IF NOT EXISTS role_permissions (
            role_id BIGINT NOT NULL,
            permission_id BIGINT NOT NULL,
            PRIMARY KEY (role_id, permission_id),
            FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
            FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
        );
        `,
		`
        CREATE TABLE IF NOT EXISTS user_roles (
            user_id CHAR(36) NOT NULL,
            role_id BIGINT NOT NULL,
            PRIMARY KEY (user_id, role_id),
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
            FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
        );
        `,
	}

	for _, query := range queries {
		if _, err := m.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration query:\n%s\nError: %w", query, err)
		}
	}
	return nil
}

func (m *MySQLStorage) CreateRole(role *models.Role) error {
	query := "INSERT INTO roles (id, name) VALUES (?, ?)"
	_, err := m.db.Exec(query, role.ID, role.Name)
	return err
}

func (m *MySQLStorage) GetRoleByID(id int64) (*models.Role, error) {
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

func (m *MySQLStorage) DeleteRole(id int64) error {
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

func (m *MySQLStorage) GetPermissionByID(id int64) (*models.Permission, error) {
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

func (m *MySQLStorage) DeletePermission(id int64) error {
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

func (m *MySQLStorage) AssignPermissionToRole(roleID, permissionID int64) error {
	query := "INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)"
	_, err := m.db.Exec(query, roleID, permissionID)
	return err
}

func (m *MySQLStorage) RemovePermissionFromRole(roleID, permissionID int64) error {
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

func (m *MySQLStorage) GetPermissionsByRole(roleID int64) ([]models.Permission, error) {
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

func (m *MySQLStorage) AssignRoleToUser(roleID int64, userID string) error {
	query := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE role_id = VALUES(role_id), updated_at = CURRENT_TIMESTAMP
	`
	_, err := m.db.Exec(query, userID, roleID)
	if err != nil {
		return fmt.Errorf("AssignRoleToUser: failed to assign role to user: %w", err)
	}
	return nil
}

func (m *MySQLStorage) RemoveRoleFromUser(roleID int64, userID string) error {
	query := `
		DELETE FROM user_roles
		WHERE user_id = ? AND role_id = ?
	`
	result, err := m.db.Exec(query, userID, roleID)
	if err != nil {
		return fmt.Errorf("RemoveRoleFromUser: failed to remove role from user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("RemoveRoleFromUser: failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("RemoveRoleFromUser: no such role assigned to user")
	}
	return nil
}

func (m *MySQLStorage) GetRoleByUser(userID string) (*models.Role, error) {
	query := `
		SELECT r.id, r.name
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ?
		LIMIT 1
	`
	row := m.db.QueryRow(query, userID)

	var role models.Role
	err := row.Scan(&role.ID, &role.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, fmt.Errorf("GetRoleByUser: failed to get role for user: %w", err)
	}
	return &role, nil
}
