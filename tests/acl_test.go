package tests

import (
	"errors"
	"sync"
	"testing"

	acllib "github.com/Volix-PriceTech/acl-lib"
	"github.com/Volix-PriceTech/acl-lib/models"
	"github.com/Volix-PriceTech/acl-lib/storage"
	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
	roles            map[int64]*models.Role
	permissions      map[int64]*models.Permission
	rolePermissions  map[int64]map[int64]bool
	userRoles        map[string]int64
	nextRoleID       int64
	nextPermissionID int64
	mu               sync.RWMutex
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		roles:            make(map[int64]*models.Role),
		permissions:      make(map[int64]*models.Permission),
		rolePermissions:  make(map[int64]map[int64]bool),
		userRoles:        make(map[string]int64),
		nextRoleID:       1,
		nextPermissionID: 1,
	}
}

func (m *MockStorage) MigrateTables() error {
	return nil
}

func (m *MockStorage) CreateRole(role *models.Role) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	role.ID = m.nextRoleID
	m.roles[role.ID] = role
	m.nextRoleID++
	return nil
}

func (m *MockStorage) GetRoleByID(id int64) (*models.Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	role, exists := m.roles[id]
	if !exists {
		return nil, storage.ErrRoleNotFound
	}
	return role, nil
}

func (m *MockStorage) DeleteRole(id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.roles, id)
	delete(m.rolePermissions, id)
	for userID, roleID := range m.userRoles {
		if roleID == id {
			delete(m.userRoles, userID)
		}
	}
	return nil
}

func (m *MockStorage) CreatePermission(permission *models.Permission) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	permission.ID = m.nextPermissionID
	m.permissions[permission.ID] = permission
	m.nextPermissionID++
	return nil
}

func (m *MockStorage) GetPermissionByID(id int64) (*models.Permission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	permission, exists := m.permissions[id]
	if !exists {
		return nil, storage.ErrPermissionNotFound
	}
	return permission, nil
}

func (m *MockStorage) DeletePermission(id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.permissions, id)
	for _, perms := range m.rolePermissions {
		delete(perms, id)
	}
	return nil
}

func (m *MockStorage) AssignPermissionToRole(roleID, permissionID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

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
	m.mu.Lock()
	defer m.mu.Unlock()

	if perms, exists := m.rolePermissions[roleID]; exists {
		delete(perms, permissionID)
		return nil
	}
	return storage.ErrPermissionNotFound
}

func (m *MockStorage) GetPermissionsByRole(roleID int64) ([]models.Permission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

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

func (m *MockStorage) AssignRoleToUser(roleID int64, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.roles[roleID]; !exists {
		return storage.ErrRoleNotFound
	}
	m.userRoles[userID] = roleID
	return nil
}

func (m *MockStorage) RemoveRoleFromUser(roleID int64, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	currentRoleID, exists := m.userRoles[userID]
	if !exists {
		return errors.New("user has no role assigned")
	}
	if currentRoleID != roleID {
		return errors.New("user is assigned a different role")
	}
	delete(m.userRoles, userID)
	return nil
}

func (m *MockStorage) GetRoleByUser(userID string) (*models.Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	roleID, exists := m.userRoles[userID]
	if !exists {
		return nil, storage.ErrRoleNotFound
	}
	role, exists := m.roles[roleID]
	if !exists {
		return nil, storage.ErrRoleNotFound
	}
	return role, nil
}

func TestACL(t *testing.T) {
	mockStorage := NewMockStorage()
	acl := acllib.NewACL(mockStorage)

	err := acl.CreateRole("admin")
	assert.NoError(t, err, "Failed to create role 'admin'")

	role, err := acl.GetRoleByID(1)
	assert.NoError(t, err, "Failed to retrieve role by ID")
	assert.Equal(t, "admin", role.Name, "Role name mismatch")

	err = acl.CreatePermission("read")
	assert.NoError(t, err, "Failed to create permission 'read'")

	permission, err := acl.GetPermissionByID(1)
	assert.NoError(t, err, "Failed to retrieve permission by ID")
	assert.Equal(t, "read", permission.Name, "Permission name mismatch")

	err = acl.AssignPermission(1, 1)
	assert.NoError(t, err, "Failed to assign permission to role")

	permissions, err := acl.GetPermissionsByRole(1)
	assert.NoError(t, err, "Failed to get permissions by role")
	assert.Len(t, permissions, 1, "Permission count mismatch")
	assert.Equal(t, "read", permissions[0].Name, "Permission name mismatch")

	err = acl.RemovePermission(1, 1)
	assert.NoError(t, err, "Failed to remove permission from role")

	permissions, err = acl.GetPermissionsByRole(1)
	assert.NoError(t, err, "Failed to get permissions by role after removal")
	assert.Len(t, permissions, 0, "Permissions should be empty after removal")

	userID := "user123"
	err = acl.AssignRoleToUser(1, userID)
	assert.NoError(t, err, "Failed to assign role to user")

	userRole, err := acl.GetRoleByUser(userID)
	assert.NoError(t, err, "Failed to get role by user")
	assert.Equal(t, "admin", userRole.Name, "User's role name mismatch")

	err = acl.RemoveRoleFromUser(1, userID)
	assert.NoError(t, err, "Failed to remove role from user")

	userRole, err = acl.GetRoleByUser(userID)
	assert.Error(t, err, "Expected error when retrieving role for user with no assigned role")
	assert.Nil(t, userRole, "User role should be nil after removal")
}

func TestAssignRoleToUser(t *testing.T) {
	mockStorage := NewMockStorage()
	acl := acllib.NewACL(mockStorage)

	err := acl.CreateRole("admin")
	assert.NoError(t, err)

	err = acl.CreateRole("user")
	assert.NoError(t, err)

	err = acl.AssignRoleToUser(1, "user1")
	assert.NoError(t, err, "Failed to assign 'admin' role to user1")

	role, err := acl.GetRoleByUser("user1")
	assert.NoError(t, err, "Failed to get role for user1")
	assert.Equal(t, "admin", role.Name, "Role name mismatch for user1")

	err = acl.AssignRoleToUser(2, "user2")
	assert.NoError(t, err, "Failed to assign 'user' role to user2")

	role, err = acl.GetRoleByUser("user2")
	assert.NoError(t, err, "Failed to get role for user2")
	assert.Equal(t, "user", role.Name, "Role name mismatch for user2")
}

func TestRemoveRoleFromUser(t *testing.T) {
	mockStorage := NewMockStorage()
	acl := acllib.NewACL(mockStorage)

	err := acl.CreateRole("admin")
	assert.NoError(t, err)

	err = acl.AssignRoleToUser(1, "user1")
	assert.NoError(t, err)

	err = acl.RemoveRoleFromUser(1, "user1")
	assert.NoError(t, err, "Failed to remove role from user1")

	role, err := acl.GetRoleByUser("user1")
	assert.Error(t, err, "Expected error when getting role for user with no assigned role")
	assert.Nil(t, role, "Role should be nil after removal")

	err = acl.RemoveRoleFromUser(1, "user1")
	assert.Error(t, err, "Expected error when removing role that is not assigned")
}

func TestGetRoleByUser(t *testing.T) {
	mockStorage := NewMockStorage()
	acl := acllib.NewACL(mockStorage)

	err := acl.CreateRole("admin")
	assert.NoError(t, err)

	err = acl.CreateRole("user")
	assert.NoError(t, err)

	err = acl.AssignRoleToUser(1, "adminUser")
	assert.NoError(t, err)

	err = acl.AssignRoleToUser(2, "regularUser")
	assert.NoError(t, err)

	adminRole, err := acl.GetRoleByUser("adminUser")
	assert.NoError(t, err, "Failed to get role for adminUser")
	assert.Equal(t, "admin", adminRole.Name, "Role name mismatch for adminUser")

	userRole, err := acl.GetRoleByUser("regularUser")
	assert.NoError(t, err, "Failed to get role for regularUser")
	assert.Equal(t, "user", userRole.Name, "Role name mismatch for regularUser")

	role, err := acl.GetRoleByUser("noRoleUser")
	assert.Error(t, err, "Expected error when getting role for user with no assigned role")
	assert.Nil(t, role, "Role should be nil for user with no assigned role")
}
