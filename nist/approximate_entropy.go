package nist

import (
	"bytes"
	"math"

	b "github.com/notJoon/drbg/bitstream"
)

func ApproximateEntropy(m uint64, bs *b.BitStream) (float64, bool, error) {
	n := uint64(bs.Len())
	var psi [2]float64

	for indexPsi := range psi {
		appendedEpsilon := make([]uint8, n+m-1)
		for i := uint64(0); i < n; i++ {
			bit, err := bs.Bit(int(i))
			if err != nil {
				return 0, false, err
			}
			appendedEpsilon[i] = bit
		}

		for j := uint64(0); j < m-1; j++ {
			bit, err := bs.Bit(int(j))
			if err != nil {
				return 0, false, err
			}
			appendedEpsilon[n+j] = bit
		}

		two_raise_power_to_m := 1 << m
		tempLastIndex := n - m + 1

		C := make([]float64, two_raise_power_to_m)
		possible_m_bits_blocks := make([][]uint8, two_raise_power_to_m)

		for bit := range possible_m_bits_blocks {
			possible_m_bits_blocks[bit] = Uint_To_BitsArray_size_N(uint64(bit), m)
			for index := uint64(0); index < tempLastIndex; index++ {
				if bytes.Equal(possible_m_bits_blocks[bit], appendedEpsilon[index:index+m]) {
					C[bit]++
				}
			}
		}

		for indexC := range C {
			C[indexC] /= float64(n)
		}

		sum := 0.0
		for _, value := range C {
			if value > 0 {
				sum += value * math.Log(value)
			}
		}
		psi[indexPsi] = sum
		m++
	}
	m--

	chi2 := 2 * float64(n) * (math.Log(2) - (psi[0] - psi[1]))
	p_val := igamc(math.Pow(2.0, float64(m-1)), chi2/2)

	return p_val, p_val > 0.01, nil
}
