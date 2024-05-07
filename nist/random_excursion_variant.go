package nist

import (
	"math"

	b "github.com/notJoon/drbg/bitstream"
)

func RandomExcursionsVariant(bs *b.BitStream) ([]float64, []bool, error) {
	n := uint64(bs.Len())
	var State_X []int64 = []int64{-9, -8, -7, -6, -5, -4, -3, -2, -1, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	var X []int64 = make([]int64, n)

	for i := uint64(0); i < n; i++ {
		bit, err := bs.Bit(int(i))
		if err != nil {
			return nil, nil, err
		}
		X[i] = 2*int64(bit) - 1
	}

	var S []int64 = make([]int64, n)
	var index_S uint64
	S[0] = X[0]
	for index_S = 1; index_S < n; index_S++ {
		S[index_S] = S[index_S-1] + X[index_S]
	}

	var S_Prime []int64 = []int64{0}
	S_Prime = append(S_Prime, S...)
	S_Prime = append(S_Prime, 0)
	S = nil

	var J int64 = 0
	for _, value := range S_Prime {
		if value == 0 {
			J++
		}
	}
	J = J - 1

	var ksi [18]int64
	for _, value := range S_Prime {
		if -9 <= value && value < 0 {
			ksi[value+9]++
		} else if 0 < value && value <= 9 {
			ksi[value+8]++
		}
	}

	var P_value []float64 = make([]float64, 18)
	var randomness []bool = make([]bool, 18)
	for i := range P_value {
		P_value[i] = math.Erfc(math.Abs(float64(ksi[i]-J)) / math.Sqrt(2.0*float64(J)*(4.0*math.Abs(float64(State_X[i]))-2.0)))
		randomness[i] = P_value[i] >= 0.01
	}

	return P_value, randomness, nil
}
