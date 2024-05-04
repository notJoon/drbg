package nist

import (
	"bytes"
	"math"

	b "github.com/notJoon/drbg/bitstream"
)

// OverlappingTemplateMatching performs the Overlapping Template Matching test from NIST SP-800-22.
// It checks for the number of occurrences of a given template B in the input bitstream.
// The bitstream is divided into N independent blocks of length M, and the test determines
// whether the number of occurrences of B in each block is approximately what would be
// expected for a random sequence.
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
func OverlappingTemplateMatching(B []uint8, eachBlockSize uint64, bs *b.BitStream) (float64, bool, error) {
	m := len(B)
	n := bs.Len()

	var (
		M uint64 = eachBlockSize   // The length of the substring of ε to be tested.
		N uint64 = (uint64(n) / M) // The number of independent blocks. N has to be fixed at 8.
	)

	// partition the sequence into N independent blocks of length M.
	blocks := make([][]uint8, N)

	// The number of occurrences of B in each block
	// by incrementing an array v[i]
	v := make([]float64, 6)

	var (
		start uint64 = 0
		end   uint64 = M
	)
	for j := range blocks {
		block := make([]uint8, M)
		for i := start; i < end; i++ {
			bit, err := bs.Bit(int(i))
			if err != nil {
				return 0, false, err
			}
			block[i-start] = bit
		}
		blocks[j] = block
		start = end
		end = end + M
	}

	// search for matches
	var numberOfOccurrence uint64
	for _, block := range blocks {
		numberOfOccurrence = 0
		for bitPos := 0; bitPos <= len(block)-m; bitPos++ {
			if bytes.Equal(block[bitPos:bitPos+m], B) {
				numberOfOccurrence++
				if numberOfOccurrence >= 5 {
					v[numberOfOccurrence]++
				}
			}
		}
	}

	// Compute values for λ, η
	_float64_m := float64(m)
	lambda := (float64(N) - _float64_m + 1) / math.Pow(2, _float64_m)
	eta := lambda / 2.0

	// Compute χ^2 as specified in Section 3.8 (p.74)
	pi := []float64{0.364091, 0.185659, 0.139381, 0.100571, 0.070432, 0.139865}

	sum := 0.0
	K := 5
	for i := 0; i < K; i++ {
		pi[i] = Pr(i, eta)
		sum += pi[i]
	}

	pi[K] = 1 - sum

	chi2 := 0.0
	_float64_N := float64(N)

	for i := range v {
		tmp := _float64_N * pi[i]
		diff := v[i] - tmp
		chi2 = diff * diff / tmp
	}

	p_value := igamc(2.5, chi2/2.0)

	return p_value, p_value >= 0.01, nil
}

// Pr calculates the probability of observing u occurrences of the template
// in a block of length M, given the expected number of occurrences (eta).
//
// Parameters:
//   - u: The number of occurrences of the template.
//   - eta: The expected number of occurrences of the template.
func Pr(u int, eta float64) float64 {
	var (
		l      int
		sum, p float64
	)

	if u == 0 {
		p = math.Exp(-1 * eta)
	} else {
		sum = 0.0
		for l = 1; l <= u; l++ {
			lgam_u, _ := math.Lgamma(float64(u))
			lgam_l, _ := math.Lgamma(float64(l))
			lgam_l_plus1, _ := math.Lgamma(float64(l + 1))
			lgam_u_l_plus1, _ := math.Lgamma(float64(u - l + 1))
			sum += math.Exp(-1*eta - float64(u)*math.Log(2) + float64(l)*math.Log(eta) - lgam_l_plus1 + lgam_u - lgam_l - lgam_u_l_plus1)
		}
		p = sum
	}
	return p
}
