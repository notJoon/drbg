package nist

import (
	"errors"
	"math"

	b "github.com/notJoon/drbg/bitstream"
)

// Runs function returns "The total number of runs" across all n bits.Runs
// The total number ofruns + the total number of one-runs.
//
// Parameters:
//   - B: The template to be searched for in the bitstream.
//
// Returns:
//   - p_value: The p-value of the test.
//   - bool: True if the test passes (p-value >= 0.01), False otherwise.
//   - error: Any error that occurred during the test, such as invalid input parameters.
func Runs(bs *b.BitStream) (float64, bool, error) {
	n := uint64(bs.Len())
	pi := 0.0

	// calculate the proportion of ones in the sequence
	for i := 0; i < int(n); i++ {
		bit, err := bs.Bit(i)
		if err != nil {
			return 0, false, err
		}
		pi += float64(bit)
	}
	pi /= float64(n)

	// determine if the prerequisite frequency test is passed
	tau := 2.0 / math.Sqrt(float64(n))
	if math.Abs(pi-0.5) >= tau {
		return 0, false, errors.New("frequency test failed")
	}

	// compute the test statistic V_n
	var prevBit byte
	V_n := 1.0
	if n > 0 {
		firstBit, _ := bs.Bit(0)
		prevBit = byte(firstBit)
	}

	for i := 1; i < int(n); i++ {
		currentBit, _ := bs.Bit(i)
		if currentBit != prevBit {
			V_n++
			prevBit = currentBit
		}
	}

	p_value := math.Erfc(math.Abs(V_n-2*float64(n)*pi*(1-pi)) / (2 * math.Sqrt(2.0*float64(n)) * pi * (1 - pi)))
	return p_value, p_value >= 0.01, nil
}
