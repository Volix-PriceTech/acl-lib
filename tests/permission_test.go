package tests

import (
	"testing"

	"github.com/Volix-PriceTech/acl-lib/models"
	"github.com/stretchr/testify/assert"
)

func TestPermission(t *testing.T) {
	permission := &models.Permission{
		ID:   1,
		Name: "edit",
	}

	assert.Equal(t, int64(1), permission.ID)
	assert.Equal(t, "edit", permission.Name)
}
