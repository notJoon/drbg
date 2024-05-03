package nist

import (
	"errors"

	b "github.com/notJoon/drbg/bitstream"
)

var (
	ErrNotEnoughLength = errors.New("input length of sequence is too small (n < 128)")
	ErrInvalidValueK   = errors.New("invalid value of K")
	ErrInvalidValueM   = errors.New("invalid value of M")
)

// LongestRunOfOnes implements the Longest Run test from NIST SP 800-22.
// The purpose of this test is to measure the longest run of ones in the input bitstream
// and evaluate the randomness of the bitstream.
//
// The test proceeds in the following steps:
//
//  1. Determine the block size (M) and the number of blocks (N) based on the length of the input bitstream.
//  2. Divide the bitstream into blocks of size M bits.
//  3. Measure the longest run of ones in each block and increase the frequency of the corresponding category.
//  4. Calculate the chi-square statistic and compute the P-value based on it.
//  5. If the P-value is greater than or equal to 0.01, the bitstream is considered random.
//
// The function returns the P-value, the test result (whether the P-value is greater than or equal to 0.01),
// and an error (if any occurred).
func LongestRunOfOnes(bs *b.BitStream) (float64, bool, error) {
	// Declare Constant
	var (
		_PI_K3_M8     = [4]float64{0.2148, 0.3672, 0.2305, 0.1875}
		_PI_K5_M128   = [6]float64{0.1174, 0.2430, 0.2493, 0.1752, 0.1027, 0.1124}
		_PI_K5_M512   = [6]float64{0.1170, 0.2460, 0.2523, 0.1755, 0.1027, 0.1124}
		_PI_K5_M1000  = [6]float64{0.1307, 0.2437, 0.2452, 0.1714, 0.1002, 0.1088}
		_PI_K6_M10000 = [7]float64{0.0882, 0.2092, 0.2483, 0.1933, 0.1208, 0.0675, 0.0727}
	)

	var M uint64 // The length of each block.
	var N uint64 // The number of blocks; selected in accordance with the value of M.
	var K uint64

	n := uint64(bs.Len())

	if n < 128 {
		return 0, false, ErrNotEnoughLength
	} else if n < 6272 {
		M = 8
		N = n / 8
		K = 3
	} else if n < 750000 {
		M = 128
		N = n / 128
		K = 5
	} else {
		M = 10000
		N = n / 10000
		K = 6
	}

	// Divide the sequence into M-bit blocks.
	sliceBoundary_start := uint64(0)
	sliceBoundary_end := M
	v := [7]uint64{0, 0, 0, 0, 0, 0, 0}
	for {
		var longest uint64 = 0
		var count uint64 = 0
		for i := sliceBoundary_start; i < sliceBoundary_end; i++ {
			bit, err := bs.Bit(int(i))
			if err != nil {
				return 0, false, err
			}
			if bit == 0 {
				longest = max(longest, count)
				count = 0
			} else {
				count++
			}
		}
		longest = max(longest, count)

		// Tabulate the frequencies νi of the longest runs of ones in each block into categories,
		// where each cell contains the number of runs of ones of a given length.
		switch K {
		case 3:
			if longest <= 1 {
				v[0]++
			} else if longest == 2 {
				v[1]++
			} else if longest == 3 {
				v[2]++
			} else {
				v[3]++
			}
		case 5:
			if longest <= 4 {
				v[0]++
			} else if longest == 5 {
				v[1]++
			} else if longest == 6 {
				v[2]++
			} else if longest == 7 {
				v[3]++
			} else if longest == 8 {
				v[4]++
			} else {
				v[5]++
			}
		case 6:
			if longest <= 10 {
				v[0]++
			} else if longest == 11 {
				v[1]++
			} else if longest == 12 {
				v[2]++
			} else if longest == 13 {
				v[3]++
			} else if longest == 14 {
				v[4]++
			} else if longest == 15 {
				v[5]++
			} else {
				v[6]++
			}
		default:
			return 0, false, ErrInvalidValueK
		}

		sliceBoundary_start += M
		sliceBoundary_end += M
		if sliceBoundary_end > n {
			break
		}
	}

	// (3) Compute Test Statistic and Reference Distribution χ^2
	var chi_square float64 = 0
	var i uint64
	switch K {
	case 3:
		for i = 0; i <= K; i++ {
			__v := float64(v[i])
			__N := float64(N)
			__PI := _PI_K3_M8[i]
			__temp := (__v - __N*__PI) * (__v - __N*__PI) / (__N * __PI)
			chi_square += __temp
		}
	case 5:
		for i = 0; i <= K; i++ {
			__v := float64(v[i])
			__N := float64(N)
			var __PI float64
			switch M {
			case 128:
				__PI = _PI_K5_M128[i]
			case 512:
				__PI = _PI_K5_M512[i]
			case 1000:
				__PI = _PI_K5_M1000[i]
			default:
				return 0, false, ErrInvalidValueM
			}
			__temp := (__v - __N*__PI) * (__v - __N*__PI) / (__N * __PI)
			chi_square += __temp
		}
	case 6:
		for i = 0; i <= K; i++ {
			__v := float64(v[i])
			__N := float64(N)
			__PI := _PI_K6_M10000[i]
			__temp := (__v - __N*__PI) * (__v - __N*__PI) / (__N * __PI)
			chi_square += __temp
		}
	}

	// (4) Compute P-value
	P_value := igamc(float64(K)/2.0, chi_square/2.0)

	return P_value, P_value >= 0.01, nil
}
