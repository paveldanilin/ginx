package slices

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoin(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}

	c := Join(a, b)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, c)
}
