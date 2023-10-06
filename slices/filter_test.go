package slices

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T) {
	r := Filter([]int{1, 2, 3, 4, 5, 6}, func(i int, t int) bool {
		return t%2 == 0
	})

	assert.Equal(t, []int{2, 4, 6}, r)
}
