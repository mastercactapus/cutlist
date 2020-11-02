package cutlist

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Length uint64

const (
	Nanometer  = Length(1)
	Millimeter = Length(1000000 * Nanometer)
	Centimeter = Length(10 * Millimeter)
	Meter      = Length(100 * Centimeter)

	Inch = Length(25400000 * Nanometer)
	Foot = Length(12 * Inch)
	Yard = Length(3 * Foot)
)

var lenRx = regexp.MustCompile(`([0-9./]+|'|"|[a-z.]+)`)

func ParseLength(s string) (Length, error) {
	s = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r
		}
		if r >= '0' && r <= '9' {
			return r
		}
		switch r {
		case '\'', '"', '.', '/':
			return r
		}

		return 0
	}, strings.ToLower(s))
	m := lenRx.FindAllString(s, -1)
	if len(m)%2 != 0 {
		return 0, fmt.Errorf("value or units missing: %s", s)
	}
	var result Length
	var val float64
	var err error
	for i := 0; i < len(m); i += 2 {
		parts := strings.SplitN(m[i], "/", 2)
		if len(parts) == 2 {
			val, err = strconv.ParseFloat(parts[0], 64)
			num := val
			if err == nil {
				val, err = strconv.ParseFloat(parts[1], 64)
			}
			if err == nil {
				val = num / val
			}
		} else {
			val, err = strconv.ParseFloat(m[i], 64)
		}
		if err != nil {
			return 0, fmt.Errorf("invalid value '%s': %w", m[i], err)
		}

		switch strings.ToLower(m[i+1]) {
		case "nm", "nanometer", "nanometers":
			result += Length(val)
		case "mm", "millimeter", "millimeters":
			result += Length(val * float64(Millimeter))
		case "cm", "centimeter", "centimeters":
			result += Length(val * float64(Centimeter))
		case "m", "meter", "meters":
			result += Length(val * float64(Meter))
		case "in", "in.", `"`, "inch", "inches":
			result += Length(val * float64(Inch))
		case "ft", "foot", "feet", "ft.", "'":
			result += Length(val * float64(Foot))
		case "yd", "yd.", "yard", "yards":
			result += Length(val * float64(Yard))
		default:
			return 0, fmt.Errorf("unknown unit '%s'", m[i+1])
		}
	}
	return result, nil
}

func (l Length) String() string {
	cm := float64(l) / float64(Centimeter)
	return fmt.Sprintf("%gcm", cm)
}
