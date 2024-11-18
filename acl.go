package acl_lib

import (
	"github.com/Volix-PriceTech/acl-lib/models"
	"github.com/Volix-PriceTech/acl-lib/storage"
)

type ACL struct {
	storage storage.Storage
}

func NewACL(s storage.Storage) *ACL {
	return &ACL{storage: s}
}

// Migrate create tables roles, permissions and roles_permissions.
func (a *ACL) Migrate() error {
	return a.storage.MigrateTables()
}

// CreateRole adds a new role with a unique ID.
func (a *ACL) CreateRole(name string) error {
	role := &models.Role{Name: name}
	return a.storage.CreateRole(role)
}

// GetRoleByID retrieves a role by its ID.
func (a *ACL) GetRoleByID(id int64) (*models.Role, error) {
	return a.storage.GetRoleByID(id)
}

// DeleteRole removes a role by its ID.
func (a *ACL) DeleteRole(id int64) error {
	return a.storage.DeleteRole(id)
}

// CreatePermission adds a new permission with a unique ID.
func (a *ACL) CreatePermission(name string) error {
	permission := &models.Permission{Name: name}
	return a.storage.CreatePermission(permission)
}

// GetPermissionByID retrieves a permission by its ID.
func (a *ACL) GetPermissionByID(id int64) (*models.Permission, error) {
	return a.storage.GetPermissionByID(id)
}

// DeletePermission removes a permission by its ID.
func (a *ACL) DeletePermission(id int64) error {
	return a.storage.DeletePermission(id)
}

// AssignPermission assigns a permission to a role.
func (a *ACL) AssignPermission(roleID, permissionID int64) error {
	return a.storage.AssignPermissionToRole(roleID, permissionID)
}

// RemovePermission removes a permission from a role.
func (a *ACL) RemovePermission(roleID, permissionID int64) error {
	return a.storage.RemovePermissionFromRole(roleID, permissionID)
}

// GetPermissionsByRole retrieves all permissions associated with a role.
func (a *ACL) GetPermissionsByRole(roleID int64) ([]models.Permission, error) {
	return a.storage.GetPermissionsByRole(roleID)
}

// AssignRoleToUser assigns a role to a user.
func (a *ACL) AssignRoleToUser(roleID int64, userID string) error {
	return a.storage.AssignRoleToUser(roleID, userID)
}

// RemoveRoleFromUser removes a role assignment from a user.
func (a *ACL) RemoveRoleFromUser(roleID int64, userID string) error {
	return a.storage.RemoveRoleFromUser(roleID, userID)
}

// GetRoleByUser retrieves the role assigned to a user.
func (a *ACL) GetRoleByUser(userID string) (*models.Role, error) {
	return a.storage.GetRoleByUser(userID)
}
