package nist

import (
	"errors"
	"fmt"
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

	zeros := 0
	ones := 0

	var S_n int64 = 0
	for i := 0; i < n; i++ {
		bit, err := bs.Bit(i)
		if err != nil {
			return 0, false, err
		}
		if bit == 0 {
			zeros += 1
			S_n -= 1
		} else {
			ones += 1
			S_n += 1
		}
	}

	fmt.Printf("Zeros: %d, Ones: %d\n", zeros, ones)

	S_obs := math.Abs(float64(S_n)) / math.Sqrt(float64(n))
	P_value := math.Erfc(S_obs / math.Sqrt(2))

	isRandom := P_value >= 0.01
	return P_value, isRandom, nil
}
