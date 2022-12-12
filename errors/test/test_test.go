package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func dummySuccessFunction() (string, error) {
	return "success", nil
}

func TestSuccess(t *testing.T) {
	result := Result(dummySuccessFunction()).Must(t)
	assert.Equal(t, "success", result)
}
