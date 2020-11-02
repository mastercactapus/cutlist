package cutlist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLength(t *testing.T) {
	check := func(exp Length, input string) {
		t.Helper()
		val, err := ParseLength(input)
		assert.NoError(t, err)
		assert.Equal(t, exp.String(), val.String())
	}

	check(1*Foot, "1ft")
	check(1*Foot+6*Inch, `1ft. 6"`)
	check(962*Millimeter, "96.2cm")
	check(127*Millimeter/10, "1/2in.")
}
