package nist

import (
	"errors"
	"fmt"
	"math"

	b "github.com/notJoon/drbg/bitstream"
)

// REPL Execute example: go run main.go -file rand_data/numbers.bin -non-overlapping-template -template "001" -block-size 100

// TODO: Receive the template string differently when calling overlapped and non-overlapped template matching tests in same call.
//
// # Non-overlapping Template Matching Test
// go run main.go -file rand_data/pcg32.bin -non-overlapping -template "000000001" -block-size 20
//
// +----------------------------------------+-------------+--------+
// | NIST STATISTICAL TEST SUITE            |     P-VALUE | RESULT |
// +----------------------------------------+-------------+--------+
// | Non-overlapping Template Matching Test |           1 | Pass   |
// +----------------------------------------+-------------+--------+
// |                                        | TOTAL TESTS | 1      |
// |                                        |        PASS | 1      |
// |                                        |        FAIL | 0      |
// +----------------------------------------+-------------+--------+
//
// # Overlapping Template Matching Test
//
// go run main.go -file rand_data/pcg32.bin -overlapping -template "111111111" -block-size 1000000
//
// +------------------------------------+-------------+--------+
// | NIST STATISTICAL TEST SUITE        |     P-VALUE | RESULT |
// +------------------------------------+-------------+--------+
// | Overlapping Template Matching Test |           1 | Pass   |
// +------------------------------------+-------------+--------+
// |                                    | TOTAL TESTS | 1      |
// |                                    |        PASS | 1      |
// |                                    |        FAIL | 0      |
// +------------------------------------+-------------+--------+

// NonOverlappingTemplateMatching performs the Non-overlapping Template Matching test from NIST SP-800-22.
// It counts the number of occurrences of a given template B in non-overlapping blocks of the input bitstream.
// The bitstream is divided into N independent blocks of length M, and the test determines whether the
// number of occurrences of B in each block is approximately what would be expected for a random sequence.
//
// Parameters:
//   - B: The template to be searched for in the bitstream.
//   - eachBlockSize: The length of each independent block (M).
//   - bs: The input bitstream.
//
// Returns:
//   - p_value: The p-value of the test.
//   - bool: True if the test passes (p-value >= 0.01), False otherwise.
//   - error: Any error that occurred during the test, such as invalid input parameters.
func NonOverlappingTemplateMatching(B []uint8, eachBlockSize uint64, bs *b.BitStream) (float64, bool, error) {
	m := len(B)
	n := bs.Len()
	M := eachBlockSize
	N := uint64(n) / M

	if uint64(n)%M != 0 {
		errorMessage := fmt.Sprintf("Input, eachBlockSize=%v, is wrong. %v mod %v remains %v", eachBlockSize, n, M, uint64(n)%M)
		return 0, false, errors.New(errorMessage)
	}

	blocks := make([][]uint8, N)
	W := make([]uint64, N)
	partitionStart := uint64(0)
	partitionEnd := M
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
		partitionEnd += M
	}

	for j := range blocks {
		for bitPosition := 0; bitPosition <= int(M)-m; bitPosition++ {
			match := true
			for i := range B {
				if blocks[j][bitPosition+i] != B[i] {
					match = false
					break
				}
			}
			if match {
				W[j]++
				bitPosition += len(B) - 1
			}
		}
	}

	_float64_m := float64(m)
	pow2m := math.Pow(2, _float64_m)
	mu := float64(M-uint64(m)+1) / pow2m
	sigma2 := float64(M) * (1/pow2m - float64(2*m-1)/math.Pow(2, 2*_float64_m))

	chi_square := 0.0
	for _, value := range W {
		chi_square = chi_square + math.Pow((float64(value)-mu), 2)/sigma2
	}

	p_value := igamc(float64(N)/2.0, chi_square/2.0)

	return p_value, p_value >= 0.01, nil
}
