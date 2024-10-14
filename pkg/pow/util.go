package pow

import (
	"math/bits"
)

const maxByte = 255

func countZeros(hashSum [32]byte) uint16 {
	var zeros uint16

	for _, b := range hashSum {
		if b == 0 {
			zeros += 8
		} else {
			zeros += uint16(bits.LeadingZeros8(b))

			return zeros
		}
	}

	if zeros > maxByte {
		return maxByte
	}

	return zeros
}
