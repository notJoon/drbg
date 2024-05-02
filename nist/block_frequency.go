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

func piWithBaseI(bs *b.BitStream, M, N uint64) ([]float64, error) {
	var sum uint64
	var ret = []float64{}
	for i := uint64(1); i <= N; i++ {
		sum = 0
		for j := uint64(0); j < M; j++ {
			bit, err := bs.Bit(int((i-1)*M + j))
			if err != nil {
				return nil, err
			}
			sum += uint64(bit)
		}
		ret = append(ret, float64(sum)/float64(M))
	}
	return ret, nil
}

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
		tempSum += (value - 0.5) * (value - 0.5)
	}
	X2 := 4 * float64(M) * tempSum

	// compute the P-value
	P_value := igamc(float64(N)/2.0, X2/2.0)

	return P_value, P_value >= 0.01, nil
}
