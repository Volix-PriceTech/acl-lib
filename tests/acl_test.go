package tests

import (
	"testing"

	"github.com/Volix-PriceTech/acl-lib"
	"github.com/Volix-PriceTech/acl-lib/models"
	"github.com/Volix-PriceTech/acl-lib/storage"
	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
	roles           map[string]*models.Role
	permissions     map[string]*models.Permission
	rolePermissions map[string]map[string]bool
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		roles:           make(map[string]*models.Role),
		permissions:     make(map[string]*models.Permission),
		rolePermissions: make(map[string]map[string]bool),
	}
}

func (m *MockStorage) CreateRole(role *models.Role) error {
	role.ID = "role_" + role.Name
	m.roles[role.ID] = role
	return nil
}

func (m *MockStorage) GetRoleByID(id string) (*models.Role, error) {
	role, exists := m.roles[id]
	if !exists {
		return nil, storage.ErrRoleNotFound
	}
	return role, nil
}

func (m *MockStorage) DeleteRole(id string) error {
	delete(m.roles, id)
	return nil
}

func (m *MockStorage) CreatePermission(permission *models.Permission) error {
	permission.ID = "perm_" + permission.Name
	m.permissions[permission.ID] = permission
	return nil
}

func (m *MockStorage) GetPermissionByID(id string) (*models.Permission, error) {
	permission, exists := m.permissions[id]
	if !exists {
		return nil, storage.ErrPermissionNotFound
	}
	return permission, nil
}

func (m *MockStorage) DeletePermission(id string) error {
	delete(m.permissions, id)
	return nil
}

func (m *MockStorage) AssignPermissionToRole(roleID, permissionID string) error {
	if _, exists := m.roles[roleID]; !exists {
		return storage.ErrRoleNotFound
	}
	if _, exists := m.permissions[permissionID]; !exists {
		return storage.ErrPermissionNotFound
	}
	if m.rolePermissions[roleID] == nil {
		m.rolePermissions[roleID] = make(map[string]bool)
	}
	m.rolePermissions[roleID][permissionID] = true
	return nil
}

func (m *MockStorage) RemovePermissionFromRole(roleID, permissionID string) error {
	if perms, exists := m.rolePermissions[roleID]; exists {
		delete(perms, permissionID)
	}
	return nil
}

func (m *MockStorage) GetPermissionsByRole(roleID string) ([]models.Permission, error) {
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

	role, err := acl.GetRoleByID("role_admin")
	assert.NoError(t, err)
	assert.Equal(t, "admin", role.Name)

	err = acl.CreatePermission("read")
	assert.NoError(t, err)

	permission, err := acl.GetPermissionByID("perm_read")
	assert.NoError(t, err)
	assert.Equal(t, "read", permission.Name)

	err = acl.AssignPermission("role_admin", "perm_read")
	assert.NoError(t, err)

	permissions, err := acl.GetPermissionsByRole("role_admin")
	assert.NoError(t, err)
	assert.Len(t, permissions, 1)
	assert.Equal(t, "read", permissions[0].Name)

	err = acl.RemovePermission("role_admin", "perm_read")
	assert.NoError(t, err)

	permissions, err = acl.GetPermissionsByRole("role_admin")
	assert.NoError(t, err)
	assert.Len(t, permissions, 0)
}
