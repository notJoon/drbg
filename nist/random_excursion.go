package nist

import (
	"math"

	b "github.com/notJoon/drbg/bitstream"
)

func RandomExcursions(bs *b.BitStream) ([]float64, []bool, error) {
	n := uint64(bs.Len())
	var State_X []int64 = []int64{-4, -3, -2, -1, 1, 2, 3, 4}

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

	var J uint64 = 0
	for _, value := range S_Prime {
		if value == 0 {
			J++
		}
	}
	J = J - 1

	Cycles := make([][]uint64, 8)
	var CycleIndex int64 = -1
	for i := range Cycles {
		Cycles[i] = make([]uint64, J)
	}
	for _, stateX := range S_Prime {
		switch stateX {
		case -4:
			Cycles[0][CycleIndex]++
		case -3:
			Cycles[1][CycleIndex]++
		case -2:
			Cycles[2][CycleIndex]++
		case -1:
			Cycles[3][CycleIndex]++
		case 0:
			CycleIndex++
		case 1:
			Cycles[4][CycleIndex]++
		case 2:
			Cycles[5][CycleIndex]++
		case 3:
			Cycles[6][CycleIndex]++
		case 4:
			Cycles[7][CycleIndex]++
		}
	}

	var v [8][6]uint64
	for rowIndex, CyclesRow := range Cycles {
		for _, occur := range CyclesRow {
			if occur < 5 {
				v[rowIndex][occur]++
			} else {
				v[rowIndex][5]++
			}
		}
	}

	chi2 := make([]float64, len(State_X))
	for chi_square_Index, x := range State_X {
		var _x float64 = float64(x)
		var pi [6]float64
		temp := 1.0 / (2.0 * math.Abs(_x))
		pi[0] = 1.0 - temp
		for k := 1; k <= 4; k++ {
			pi[k] = temp * temp * math.Pow((1.0-temp), float64(k-1))
		}
		pi[5] = temp * math.Pow((1-temp), 4)

		var sum float64 = 0.0
		var J_pi float64
		for k := 0; k <= 5; k++ {
			J_pi = float64(J) * pi[k]
			sum += (float64(v[chi_square_Index][k]) - J_pi) * (float64(v[chi_square_Index][k]) - J_pi) / J_pi
		}
		chi2[chi_square_Index] = sum
	}

	p_value := make([]float64, 8)
	randomness := make([]bool, 8)

	for i := range p_value {
		p_value[i] = igamc(5.0/2.0, chi2[i]/2.0)
		randomness[i] = p_value[i] >= 0.01
	}

	return p_value, randomness, nil
}
