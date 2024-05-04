package nist

import (
	"errors"
	"math"

	b "github.com/notJoon/drbg/bitstream"
)

var (
	ErrEmptyBitStream = errors.New("empty bitstream")
)

// FrequencyTest evaluates the randomness of a sequence of numbers ranging from 0 to 255,
// which are read from a file and used as a bitstream. This function applies the monobit test,
// a basic test of randomness, which counts the number of 0s and 1s in the bitstream and
// assesses whether their distribution is close enough to 50/50 as would be expected in a
// random sequence.
//
// The test calculates the test statistic S_n as follows:
//
//	S_n = |sum from i=1 to n of (2 * x_i - 1)|
//
// where x_i represents each bit in the sequence (0 or 1).
//
// This statistic is then normalized to account for the length of the sequence:
//
//	S_obs = S_n / sqrt(n)
//
// where n is the total number of bits.
//
// The final result of the test is determined by calculating a P-value using the complementary
// error function:
//
//	P = erfc(S_obs / sqrt(2))
//
// A P-value of 0.01 or higher suggests that the sequence can be considered
// random with a confidence level of 99%. If the P-value is less than 0.01, the sequence is
// considered non-random, indicating potential patterns or biases in the distribution of bits.
//
// Returns the P-value of the test, a boolean indicating whether the sequence is considered
// random, and an error if there is an issue with the input bitstream.
//
// Parameters:
//   - B: The template to be searched for in the bitstream.
//
// Returns:
//   - p_value: The p-value of the test.
//   - bool: True if the test passes (p-value >= 0.01), False otherwise.
//   - error: Any error that occurred during the test, such as invalid input parameters.
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
	p_value := math.Erfc(S_obs / math.Sqrt2)

	isRandom := p_value >= 0.01
	return p_value, isRandom, nil
}
