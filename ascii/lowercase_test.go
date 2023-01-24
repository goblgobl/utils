package ascii

import (
	"strings"
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Lowercase(t *testing.T) {
	largeUp := strings.Repeat("A", 2000)
	largeDown := strings.ToLower(largeUp)

	cases := map[string]string{
		"over":        "over",
		"OVER":        "over",
		"Over 9000!":  "over 9000!",
		"OvER__9000!": "over__9000!",
		largeUp:       largeDown,
		largeDown:     largeDown,
	}

	for input, expected := range cases {
		assert.Equal(t, Lowercase(input), expected)
	}
}
