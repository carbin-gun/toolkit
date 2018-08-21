package base36

import (
	"encoding/binary"
	"crypto/rand"
	"strings"
)

const (
	BaseOf36      uint64 = 36
	CharacterSize        = 16
	ResultSize           = CharacterSize + 3
)

var (
	mapping = []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7',
		'8', '9'}
)

func Rand() string {
	var size = 0
	var number uint64
	var tmp string
	var increment = 0
	var buf = new(strings.Builder)
	buf.Grow(CharacterSize)

	for size < 16 {
		binary.Read(rand.Reader, binary.LittleEndian, &number)
		tmp = uint64ToBase36(number)
		if size+len(tmp) >= CharacterSize {
			increment = CharacterSize - size
		} else {
			increment = len(tmp)
		}
		size += increment
		buf.WriteString(tmp[0:increment])

	}
	return format(buf)
}

func format(src *strings.Builder) string {
	var base = src.String()
	var builder = new(strings.Builder)
	builder.Grow(ResultSize)
	builder.WriteString(base[0:4])
	builder.WriteRune('-')

	builder.WriteString(base[4:8])
	builder.WriteRune('-')

	builder.WriteString(base[8:12])
	builder.WriteRune('-')

	builder.WriteString(base[12:16])

	var result = builder.String()
	return result
}
func uint64ToBase36(u uint64) string {
	var builder strings.Builder
	var base = u
	var reminder uint64 = 0
	for base != 0 {
		reminder = base % BaseOf36
		base = base / BaseOf36
		builder.WriteRune(mapping[reminder])
	}
	return builder.String()
}
