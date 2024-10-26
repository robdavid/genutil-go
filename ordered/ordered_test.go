package ordered

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	assert.Equal(t, 1, Min(4, 3, 2, 1))
}

func TestMax(t *testing.T) {
	assert.Equal(t, 4, Max(1, 2, 3, 4, 0))
}

func TestMinString(t *testing.T) {
	assert.Equal(t, "lamb", Min("mary", "little", "lamb"))
}

func TestMaxString(t *testing.T) {
	assert.Equal(t, "mary", Max("mary", "little", "lamb"))
}
