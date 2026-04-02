package position

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosition_Fields(t *testing.T) {
	pos := Position{
		ID:        1,
		CompanyID: 1,
		Name:      "Engineer",
		Status:    1,
	}

	assert.Equal(t, int64(1), pos.ID)
	assert.Equal(t, int64(1), pos.CompanyID)
	assert.Equal(t, "Engineer", pos.Name)
	assert.Equal(t, int8(1), pos.Status)
}

func TestPosition_DefaultStatus(t *testing.T) {
	pos := Position{
		Name: "Manager",
	}
	// Status defaults to 0
	assert.Equal(t, int8(0), pos.Status)
}

func TestPosition_WithParentID(t *testing.T) {
	pos := Position{
		ID:       1,
		ParentID: ptrInt64(10),
		Name:     "Senior Engineer",
		Status:   1,
	}

	assert.NotNil(t, pos.ParentID)
	assert.Equal(t, int64(10), *pos.ParentID)
}

func ptrInt64(v int64) *int64 {
	return &v
}
