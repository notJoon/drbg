package nist

import (
	"errors"
	"fmt"
	"math"

	b "github.com/notJoon/drbg/bitstream"
)

// REPL Execute example: go run main.go -file rand_data/numbers.bin -non-overlapping-template -template "001" -block-size 100

func NonOverlappingTemplateMatching(B []uint8, eachBlockSize uint64, bs *b.BitStream) (float64, bool, error) {
	m := len(B)
	n := bs.Len()

	var (
		M uint64 = eachBlockSize
		N uint64
	)

	if uint64(n)%M != 0 {
		errorMessage := fmt.Sprintf("Input, eachBlockSize=%v, is wrong. %v mod %v remains %v", eachBlockSize, n, M, uint64(n)%M)
		return 0, false, errors.New(errorMessage)
	}
	N = (uint64(n) / M)

	blocks := make([][]uint8, N)
	W := make([]uint64, N)
	partitionStart := uint64(0)
	partitionEnd := uint64(M)
	for j := range blocks {
		block := make([]uint8, M)
		for i := partitionStart; i < partitionEnd; i++ {
			bit, err := bs.Bit(int(i))
			if err != nil {
				return 0, false, err
			}
			block[i-partitionStart] = bit
		}
		blocks[j] = block
		partitionStart = partitionEnd
		partitionEnd = partitionEnd + M
	}

	for j := range blocks {
		for bitPosition := 0; bitPosition <= int(M)-m; bitPosition++ {
			for i := range B {
				if blocks[j][bitPosition+i] != B[i] {
					goto UN_HIT
				}
			}
			W[j]++
			bitPosition = bitPosition + len(B) - 1
		UN_HIT:
			// Do nothing
		}
	}

	var mu, sigma2 float64
	_float64_m := float64(m)
	mu = float64(M-uint64(m)+1) / math.Pow(2, _float64_m)
	sigma2 = float64(M) * (1/math.Pow(2, _float64_m) - float64(2*m-1)/math.Pow(2, 2*_float64_m))

	chi_square := 0.0
	for _, value := range W {
		chi_square = chi_square + math.Pow((float64(value)-mu), 2)/sigma2
	}

	p_value := igamc(float64(N)/2.0, chi_square/2.0)

	return p_value, p_value >= 0.01, nil
}
