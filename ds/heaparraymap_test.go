package ds

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeapMapRemove(t *testing.T) {
	var (
		h          HeapArrayMap[int, int, int]
		i, k, v, s int
	)
	h.Push(3, 3, 5)
	h.Push(1, 1, 2)
	h.Push(2, 2, 4)
	h.Push(4, 4, 1)
	h.Push(5, 5, 3)

	assert.True(t, h.check())

	i, v = h.Get(1)
	assert.Equal(t, 1, i)
	assert.Equal(t, 1, v)

	i, v = h.Get(2)
	assert.Equal(t, 2, i)
	assert.Equal(t, 2, v)
	h.Remove(i)
	assert.True(t, h.check())

	i, k, v, s = h.Top()
	assert.Equal(t, 3, i)
	assert.Equal(t, 4, k)
	assert.Equal(t, 4, v)
	assert.Equal(t, 1, s)
	h.Pop()

	assert.True(t, h.check())
}

func TestFuzzy(t *testing.T) {
	n := 10000
	m := 1000
	var h HeapArrayMap[int, int, int]
	for b := 0; b < n; b++ {
		k := rand.N(m) + 1
		i, v := h.Get(k)
		s := rand.N(m) + 1
		if i >= 0 {
			assert.Equal(t, k, v)
			if rand.N(2) == 0 {
				h.Update(i, s)
			} else {
				h.Remove(i)
			}
		} else {
			h.Push(k, k, s)
		}
		assert.True(t, h.check())
	}
	fmt.Println(h.Size())
}
