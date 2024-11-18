package storage

import (
	"errors"

	"github.com/Volix-PriceTech/acl-lib/models"
)

var ErrRoleNotFound = errors.New("role not found")
var ErrPermissionNotFound = errors.New("permission not found")

type Storage interface {
	CreateRole(role *models.Role) error
	GetRoleByID(id int64) (*models.Role, error)
	DeleteRole(id int64) error
	CreatePermission(permission *models.Permission) error
	GetPermissionByID(id int64) (*models.Permission, error)
	DeletePermission(id int64) error
	AssignPermissionToRole(roleID, permissionID int64) error
	RemovePermissionFromRole(roleID, permissionID int64) error
	GetPermissionsByRole(roleID int64) ([]models.Permission, error)
	AssignRoleToUser(roleID int64, userID string) error
	RemoveRoleFromUser(roleID int64, userID string) error
	GetRoleByUser(userID string) (*models.Role, error)
	MigrateTables() error
}
