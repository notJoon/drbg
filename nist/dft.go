package nist

import (
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
	b "github.com/notJoon/drbg/bitstream"
)

// ref: A Statistical Test Suite for Random and Pseudorandom Number Generators for Cryptographic Application
// Section 2.6 Discrete Fourier Transform (Spectral) Test

// DFT performs the Discrete Fourier Transform (Spectral) test on the given bitstream.
// The purpose of this test is to detect periodic features in the input sequence that would indicate a deviation from the assumption of randomness.
// The test uses the discrete Fourier transform to calculate the magnitude of the Fourier coefficients of the input sequence.
// It then compares the peak heights (the moduli of the Fourier coefficients) to a threshold value.
// If the number of peaks exceeding the threshold is significantly different from the expected number (95% of n/2),
// the sequence is considered non-random.
func DFT(bs *b.BitStream) (float64, bool, error) {
	n := bs.Len()
	X := make([]float64, 0, n/8)

	// Bitwise processing when generating X.
	// Processing in bytes can cause the bit order to be reversed.
	for i := 0; i < n; i++ {
		bit, err := bs.Bit(i)
		if err != nil {
			return 0, false, err
		}
		X = append(X, 2*float64(bit)-1)
	}

	// apply DFT on X to produce S := DFT(X)
	S := fft.FFTReal(X)

	// calculate M = modulus(S´) ≡ |S'|
	// Where S' is the substring consisting of the first n/2 elements of S
	// and the modulus function produces a sequence of peak heights.
	M := modulus(S)

	// Compute T = sqrt(log(1/0.05) * n) => 95% peak height threshold value.
	// Under an assumption of randomness, 95% of the value obtained from the test should be less than T.
	const threshold = 2.995732274 // log(1/0.05) ≈ 2.995732274
	T := math.Sqrt(threshold * float64(n))

	// Compute N_0 = 0.95n/2
	// N_0 is the expected theorertical (95%) number of peaks that are less than T.
	N0 := 0.95 * float64(n) / 2

	// N_1 is the actual observed number of peaks in M that are less than T.
	N1 := 0
	for _, value := range M {
		if value < T {
			N1++
		}
	}

	d := (float64(N1) - N0) / math.Sqrt(float64(n)*0.95*0.05/4)

	p_value := math.Erfc(math.Abs(d) / math.Sqrt2)

	return p_value, p_value >= 0.01, nil
}

func modulus(input []complex128) []float64 {
	half := len(input) / 2
	result := make([]float64, half)

	for i := 0; i < half; i++ {
		result[i] = cmplx.Abs(input[i])
	}

	return result
}
