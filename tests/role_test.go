package tests

import (
	"testing"

	"github.com/Volix-PriceTech/acl-lib/models"
	"github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {
	role := &models.Role{
		ID:   1,
		Name: "user",
	}

	assert.Equal(t, int64(1), role.ID)
	assert.Equal(t, "user", role.Name)
}
