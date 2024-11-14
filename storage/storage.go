package storage

import (
	"errors"

	"github.com/Volix-PriceTech/acl-lib/models"
)

var ErrRoleNotFound = errors.New("role not found")
var ErrPermissionNotFound = errors.New("permission not found")

type Storage interface {
	CreateRole(role *models.Role) error
	GetRoleByID(id string) (*models.Role, error)
	DeleteRole(id string) error
	CreatePermission(permission *models.Permission) error
	GetPermissionByID(id string) (*models.Permission, error)
	DeletePermission(id string) error
	AssignPermissionToRole(roleID, permissionID string) error
	RemovePermissionFromRole(roleID, permissionID string) error
	GetPermissionsByRole(roleID string) ([]models.Permission, error)
	MigrateTables() error
}
