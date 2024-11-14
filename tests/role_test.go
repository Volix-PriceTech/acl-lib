package tests

import (
	"testing"

	"github.com/Volix-PriceTech/acl-lib/models"
	"github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {
	role := &models.Role{
		ID:   "role_user",
		Name: "user",
	}

	assert.Equal(t, "role_user", role.ID)
	assert.Equal(t, "user", role.Name)
}
