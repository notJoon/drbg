package nist

// ref: https://gist.github.com/StuartGordonReid/821b002a307ce3757d27

import (
	"errors"
	"fmt"

	b "github.com/notJoon/drbg/bitstream"
)

var ErrSequenceTooShort = errors.New("input sequence length should be at least 100 bits")

// piWithBaseI calculates the proportion πi of '1's in each M-bit block of the bitstream.
// It returns a slice of float64 representing the proportion of '1's for each block
// and an error if an issue occurs during bit extraction.
func piWithBaseI(bs *b.BitStream, M, N uint64) ([]float64, error) {
	var sum uint64
	result := make([]float64, 0, N)

	for i := uint64(1); i <= N; i++ {
		sum = 0
		for j := uint64(0); j < M; j++ {
			bit, err := bs.Bit(int((i-1)*M + j))
			if err != nil {
				return nil, err
			}
			sum += uint64(bit)
		}
		result = append(result, float64(sum)/float64(M))
	}
	return result, nil
}

// BlockFrequencyTest performs the Frequency Test Within a Block as defined in NIST SP800-22.
// It takes a BitStream and block size M as input and returns the P-value of the test,
// a bool representing if the P-value suggests randomness (true if P >= 0.01), and an error if any.
//
// Parameters:
//   - B: The template to be searched for in the bitstream.
//
// Returns:
//   - p_value: The p-value of the test.
//   - bool: True if the test passes (p-value >= 0.01), False otherwise.
//   - error: Any error that occurred during the test, such as invalid input parameters.
func BlockFrequencyTest(bs *b.BitStream, M uint64) (float64, bool, error) {
	n := uint64(bs.Len())
	if n < 100 {
		return 0, false, fmt.Errorf("input sequence length should be at least 100 bits, got %d", n)
	}
	if M < 20 || M <= n/100 {
		maxM := n / 100
		return 0, false, fmt.Errorf("invalid block size. got %d, should be at least 20 and less than %d", M, maxM)
	}

	// partition the input sequence into N = floor(n/M) non-overlapping blocks
	N := n / M

	// determine the proportion πi of ones in each M-bit block
	pi, err := piWithBaseI(bs, M, N)
	if err != nil {
		return 0, false, err
	}

	// compute the test statistic X^2
	tempSum := 0.0
	for _, value := range pi {
		diff := value - 0.5
		tempSum += diff * diff
	}
	X2 := 4 * float64(M) * tempSum

	// compute the P-value using the incomplete gamma function complement
	p_value := igamc(float64(N)/2.0, X2/2.0)

	return p_value, p_value >= 0.01, nil
}
