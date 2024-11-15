package tests

import (
	"testing"

	"github.com/Volix-PriceTech/acl-lib"
	"github.com/Volix-PriceTech/acl-lib/models"
	"github.com/Volix-PriceTech/acl-lib/storage"
	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
	roles           map[int64]*models.Role
	permissions     map[int64]*models.Permission
	rolePermissions map[int64]map[int64]bool
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		roles:           make(map[int64]*models.Role),
		permissions:     make(map[int64]*models.Permission),
		rolePermissions: make(map[int64]map[int64]bool),
	}
}

func (m *MockStorage) MigrateTables() error {
	return nil
}

func (m *MockStorage) CreateRole(role *models.Role) error {
	role.ID = 1
	m.roles[role.ID] = role
	return nil
}

func (m *MockStorage) GetRoleByID(id int64) (*models.Role, error) {
	role, exists := m.roles[id]
	if !exists {
		return nil, storage.ErrRoleNotFound
	}
	return role, nil
}

func (m *MockStorage) DeleteRole(id int64) error {
	delete(m.roles, id)
	return nil
}

func (m *MockStorage) CreatePermission(permission *models.Permission) error {
	permission.ID = 1
	m.permissions[permission.ID] = permission
	return nil
}

func (m *MockStorage) GetPermissionByID(id int64) (*models.Permission, error) {
	permission, exists := m.permissions[id]
	if !exists {
		return nil, storage.ErrPermissionNotFound
	}
	return permission, nil
}

func (m *MockStorage) DeletePermission(id int64) error {
	delete(m.permissions, id)
	return nil
}

func (m *MockStorage) AssignPermissionToRole(roleID, permissionID int64) error {
	if _, exists := m.roles[roleID]; !exists {
		return storage.ErrRoleNotFound
	}
	if _, exists := m.permissions[permissionID]; !exists {
		return storage.ErrPermissionNotFound
	}
	if m.rolePermissions[roleID] == nil {
		m.rolePermissions[roleID] = make(map[int64]bool)
	}
	m.rolePermissions[roleID][permissionID] = true
	return nil
}

func (m *MockStorage) RemovePermissionFromRole(roleID, permissionID int64) error {
	if perms, exists := m.rolePermissions[roleID]; exists {
		delete(perms, permissionID)
	}
	return nil
}

func (m *MockStorage) GetPermissionsByRole(roleID int64) ([]models.Permission, error) {
	var perms []models.Permission
	if permissionIDs, exists := m.rolePermissions[roleID]; exists {
		for permID := range permissionIDs {
			if perm, exists := m.permissions[permID]; exists {
				perms = append(perms, *perm)
			}
		}
	}
	return perms, nil
}

func TestACL(t *testing.T) {
	mockStorage := NewMockStorage()
	acl := acl_lib.NewACL(mockStorage)

	err := acl.CreateRole("admin")
	assert.NoError(t, err)

	role, err := acl.GetRoleByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "admin", role.Name)

	err = acl.CreatePermission("read")
	assert.NoError(t, err)

	permission, err := acl.GetPermissionByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "read", permission.Name)

	err = acl.AssignPermission(1, 1)
	assert.NoError(t, err)

	permissions, err := acl.GetPermissionsByRole(1)
	assert.NoError(t, err)
	assert.Len(t, permissions, 1)
	assert.Equal(t, "read", permissions[0].Name)

	err = acl.RemovePermission(1, 1)
	assert.NoError(t, err)

	permissions, err = acl.GetPermissionsByRole(1)
	assert.NoError(t, err)
	assert.Len(t, permissions, 0)
}
