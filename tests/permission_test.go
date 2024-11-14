package tests

import (
	"testing"

	"github.com/Volix-PriceTech/acl-lib/models"
	"github.com/stretchr/testify/assert"
)

func TestPermission(t *testing.T) {
	permission := &models.Permission{
		ID:   "perm_edit",
		Name: "edit",
	}

	assert.Equal(t, "perm_edit", permission.ID)
	assert.Equal(t, "edit", permission.Name)
}
