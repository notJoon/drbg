package nist

import (
	"fmt"
	"math"

	b "github.com/notJoon/drbg/bitstream"
)

func CumulativeSums(mode int, bs *b.BitStream) (float64, bool, error) {
	n := uint64(bs.Len())

	if n < 2 {
		panic(fmt.Sprintf("input length is too short, should be larger than 2. got=%d", n))
	}

	X := make([]int8, n)
	S := make([]int64, n)

	for i := uint64(0); i < n; i++ {
		bit, err := bs.Bit(int(i))
		if err != nil {
			return 0, false, err
		}

		X[i] = 2*int8(bit) - 1
	}

	switch mode {
	case 0:
		_n := int64(n)
		index_s := int64(0)

		S[index_s] = int64(X[index_s])
		for index_s = 1; index_s < _n; index_s++ {
			S[index_s] = S[index_s-1] + int64(X[index_s])
		}
	case 1:
		index_x := uint64(n - 1)
		index_s := uint64(0)

		S[index_s] = int64(X[index_x])
		for index_s = 1; index_s < n; index_s++ {
			index_x--
			S[index_s] = S[index_s-1] + int64(X[index_x])
		}
	default:
		panic(fmt.Sprintf("invalid mode: %d", mode))
	}

	z := math.Abs(float64(S[0]))
	now := 0.0

	for index := uint64(1); index < n; index++ {
		now = math.Abs(float64(S[index]))
		if now > z {
			z = now
		}
	}

	n_float64 := float64(n)

	var k int64
	term1, term2 := 0.0, 0.0

	sqrt_n := math.Sqrt(n_float64)

	for k = int64((-1.0*n_float64/z + 1.0) / 4.0); k <= int64((n_float64/z-1.0)/4.0); k++ {
		term1 += cumulativeDistibution(float64(4*k+1)*z/sqrt_n) - cumulativeDistibution(float64(4*k-1)*z/sqrt_n)
	}

	for k = int64((-1.0*n_float64/z - 3.0) / 4.0); k <= int64((n_float64/z-1.0)/4.0); k++ {
		term2 += cumulativeDistibution(float64(4*k+3)*z/sqrt_n) - cumulativeDistibution(float64(4*k+1)*z/sqrt_n)
	}

	p_value := 1.0 - term1 + term2

	return p_value, p_value >= 0.01, nil
}

func cumulativeDistibution(z float64) float64 {
	return 0.5 * (math.Erf(z/math.Sqrt2) + 1)
}
