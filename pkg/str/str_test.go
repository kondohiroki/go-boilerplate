package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	list := []string{"apple", "banana", "orange"}
	assert.True(t, StringInSlice("banana", list))
	assert.False(t, StringInSlice("pineapple", list))
}
