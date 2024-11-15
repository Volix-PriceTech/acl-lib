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

func (a *ACL) Migrate() error {
	return a.storage.MigrateTables()
}

func (a *ACL) CreateRole(name string) error {
	role := &models.Role{Name: name}
	return a.storage.CreateRole(role)
}

func (a *ACL) GetRoleByID(id int64) (*models.Role, error) {
	return a.storage.GetRoleByID(id)
}

func (a *ACL) DeleteRole(id int64) error {
	return a.storage.DeleteRole(id)
}

func (a *ACL) CreatePermission(name string) error {
	permission := &models.Permission{Name: name}
	return a.storage.CreatePermission(permission)
}

func (a *ACL) GetPermissionByID(id int64) (*models.Permission, error) {
	return a.storage.GetPermissionByID(id)
}

func (a *ACL) DeletePermission(id int64) error {
	return a.storage.DeletePermission(id)
}

func (a *ACL) AssignPermission(roleID, permissionID int64) error {
	return a.storage.AssignPermissionToRole(roleID, permissionID)
}

func (a *ACL) RemovePermission(roleID, permissionID int64) error {
	return a.storage.RemovePermissionFromRole(roleID, permissionID)
}

func (a *ACL) GetPermissionsByRole(roleID int64) ([]models.Permission, error) {
	return a.storage.GetPermissionsByRole(roleID)
}
