package ascii

import (
	"math/rand"
	"strconv"
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Atoi(t *testing.T) {
	for i := 0; i < 1000; i++ {
		n := rand.Int()

		s1 := strconv.Itoa(n)
		n1, e1 := Atoi(s1)
		assert.Equal(t, n, n1)
		assert.Equal(t, e1, "")

		s2 := s1 + "-suffix"
		n2, e2 := Atoi(s2)
		assert.Equal(t, n, n2)
		assert.Equal(t, e2, "-suffix")
	}
}
