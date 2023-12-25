package ascii

import (
	"math"
	"math/rand"
	"strconv"
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Atof(t *testing.T) {
	f, s := Atof("123.456 Hi")
	assert.Equal(t, f, 123.456)
	assert.Equal(t, s, " Hi")

	f, s = Atof("2")
	assert.Equal(t, f, 2)
	assert.Equal(t, s, "")

	for i := 0; i < 100; i++ {
		actual := rand.Float64() * math.Pow10(rand.Intn(10))
		f, s = Atof(strconv.FormatFloat(actual, 'f', 15, 64))
		assert.Delta(t, f, actual, 0.00001)
		assert.Equal(t, s, "")
	}
}
