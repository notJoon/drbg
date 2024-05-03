package nist

// ref: https://gist.github.com/StuartGordonReid/821b002a307ce3757d27

import (
	"errors"

	b "github.com/notJoon/drbg/bitstream"
)

var (
	ErrSequenceTooShort = errors.New("input sequence length should be at least 100 bits")
	ErrInvalidBlockSize = errors.New("invalid block size M")
)

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
func BlockFrequencyTest(bs *b.BitStream, M uint64) (float64, bool, error) {
	n := uint64(bs.Len())
	if n < 100 {
		return 0, false, ErrSequenceTooShort
	}
	if M < 20 || M <= n/100 {
		return 0, false, ErrInvalidBlockSize
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
	P_value := igamc(float64(N)/2.0, X2/2.0)

	return P_value, P_value >= 0.01, nil
}
