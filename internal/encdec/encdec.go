package encdec

import (
	"hash/fnv"
	"math"
	"strings"
)

const (
	chars  = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0123456789_"
	length = len(chars)
)

func Encode(url string) (string, uint64) {
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(10)
	h := fnv.New64a()
	h.Write([]byte(url))
	number := h.Sum64()
	var key uint64
	for ; number > 0; number = number / uint64(length) {
		if len(encodedBuilder.String()) == 10 {
			break
		}
		encodedBuilder.WriteByte(chars[(number % uint64(length))])
	}

	key = Decode(encodedBuilder.String())

	return encodedBuilder.String(), key
}

func Decode(encoded string) uint64 {
	var number uint64

	for i, symbol := range encoded {
		alphabeticPosition := strings.IndexRune(chars, symbol)

		if alphabeticPosition == -1 {
			return uint64(alphabeticPosition)
		}
		number += uint64(alphabeticPosition) * uint64(math.Pow(float64(length), float64(i)))
	}

	return number
}
