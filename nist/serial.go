package nist

import (
	"bytes"
	"math"

	b "github.com/notJoon/drbg/bitstream"
)

func Serial(m uint64, bs *b.BitStream) ([]float64, []bool, error) {
	n := uint64(bs.Len())

	v := make([][]uint64, 3)

	var section2_idx uint64
	for section2_idx = 0; section2_idx <= 2; section2_idx++ {
		if int64(m)-int64(section2_idx)-1 <= 0 {
			break
		}

		appendedEpsilon := make([]uint8, n+m-section2_idx-1)
		for i := uint64(0); i < n; i++ {
			bit, err := bs.Bit(int(i))
			if err != nil {
				return nil, nil, err
			}
			appendedEpsilon[i] = bit
		}
		for i := uint64(0); i < m-section2_idx-1; i++ {
			bit, err := bs.Bit(int(i))
			if err != nil {
				return nil, nil, err
			}
			appendedEpsilon[n+i] = bit
		}
		blockSize := m - section2_idx
		v[section2_idx] = make([]uint64, uint64(math.Pow(2, float64(blockSize))))

		// determine the frequency of all possible overlapping m-bit blocks
		// the frequency of all possible overlapping m-bit blocks.
		for blockIndex := uint64(0); blockIndex < n; blockIndex++ {
			for vIndex := range v[section2_idx] {
				arr1 := appendedEpsilon[blockIndex : blockIndex+blockSize]
				arr2 := Uint_To_BitsArray_size_N(uint64(vIndex), blockSize)
				if bytes.Equal(arr1, arr2) {
					v[section2_idx][vIndex]++
				}
			}
		}
	}

	// compute ψ
	// // ψ_m = psi[0], ψ_{m-1} = psi[1], ψ_{m-2} = psi[2]
	psi := [3]float64{0, 0, 0}
	for i := range psi {
		if len(v[i]) == 0 {
			break
		}
		for _, value := range v[i] {
			psi[i] += math.Pow(float64(value), 2)
		}
		psi[i] = math.Pow(2, float64(m)-float64(i))/float64(n)*psi[i] - float64(n)
	}

	// Compute ∇ψ^2 and ∇^2ψ^2
	delta1 := psi[0] - psi[1]
	delta2 := psi[0] - 2*psi[1] + psi[2]

	temp := math.Pow(2, float64(m-2))
	p1 := igamc(temp, delta1/2)
	p2 := igamc(temp/2, delta2/2)

	p_val := []float64{p1, p2}
	pass := []bool{p1 >= 0.01, p2 >= 0.01}

	return p_val, pass, nil
}

func Uint_To_BitsArray_size_N(input uint64, N uint64) (bitArray []uint8) {
	bitArray = make([]uint8, N)
	if input == 0 {
		return
	}

	var index uint64 = N - 1
	var quotient, remainer uint64
	for {
		quotient = input / 2
		remainer = input - quotient*2
		bitArray[index] = uint8(remainer)
		index--
		input = quotient
		if quotient == 0 {
			return
		}
	}
}
