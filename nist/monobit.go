package nist

import (
	"errors"
	"math"

	b "github.com/notJoon/drbg/bitstream"
)

var (
	ErrEmptyBitStream = errors.New("empty bitstream")
)

func FrequencyTest(bs *b.BitStream) (float64, bool, error) {
	n := bs.Len()
	if n == 0 {
		return 0, false, ErrEmptyBitStream
	}

	var S_n int64 = 0
	for i := 0; i < n; i++ {
		bit, err := bs.Bit(i)
		if err != nil {
			return 0, false, err
		}
		if bit == 0 {
			S_n -= 1
		} else {
			S_n += 1
		}
	}

	S_obs := math.Abs(float64(S_n)) / math.Sqrt(float64(n))
	P_value := math.Erfc(S_obs / math.Sqrt(2))

	isRandom := P_value >= 0.01
	return P_value, isRandom, nil
}
