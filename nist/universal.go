package nist

import (
	"math"

	b "github.com/notJoon/drbg/bitstream"
)

func UniversalRecommendedValues(bs *b.BitStream) (float64, bool, error) {
	n := uint64(bs.Len())
	L, Q := recommandedInputSize(n)
	return Universal(L, Q, n, bs)
}

func recommandedInputSize(n uint64) (L uint64, Q uint64) {
	if n >= 1059061760 {
		L = 16
		Q = 655360
	} else if n >= 496435200 {
		L = 15
		Q = 327680
	} else if n >= 231669760 {
		L = 14
		Q = 163840
	} else if n >= 107560960 {
		L = 13
		Q = 81920
	} else if n >= 49643520 {
		L = 12
		Q = 40960
	} else if n >= 22753280 {
		L = 11
		Q = 20480
	} else if n >= 10342400 {
		L = 10
		Q = 10240
	} else if n >= 4654080 {
		L = 9
		Q = 5120
	} else if n >= 2068480 {
		L = 8
		Q = 2560
	} else if n >= 904960 {
		L = 7
		Q = 1280
	} else if n >= 387840 {
		L = 6
		Q = 640
	} else {
		panic("minimum input size is 387840 bits")
	}
	return
}

func array2BinaryInt(arr []uint8) uint64 {
	digits := uint64(1)
	index_T := uint64(0)

	// divide into L-bits
	for i := len(arr) - 1; i >= 0; i-- {
		index_T += uint64(arr[i]) * digits
		digits *= 2
	}

	return index_T
}

// input size recommendation
// n >= (Q + K) * L
// 6 <= L <= 16, Q = 10 * 2^L, k = floor(n/L) - Q ~= 1000 * 2^L
func Universal(L uint64, Q uint64, n uint64, bs *b.BitStream) (float64, bool, error) {
	expectedValue_mu := [16]float64{0.7326495, 1.5374383, 2.4016068, 3.3112247, 4.2534266, 5.2177052, 6.1962507, 7.1836656, 8.1764248, 9.1723243, 10.170032, 11.168765, 12.168070, 13.167693, 14.167488, 15.167379}
	variance_sigma := [16]float64{0.690, 1.338, 1.901, 2.358, 2.705, 2.954, 3.125, 3.238, 3.311, 3.356, 3.384, 3.401, 3.410, 3.416, 3.419, 3.421}

	K := (n / L) - Q
	_float64_Q := float64(Q)

	blocks := make([][]uint8, 0, Q+K)
	T := make([]float64, Q)

	// Divide into L-bits
	var blockNum uint64 = 0
	for {
		block := make([]uint8, L)
		for i := uint64(0); i < L; i++ {
			bit, err := bs.Bit(int(blockNum*L + i))
			if err != nil {
				return 0, false, err
			}
			block[i] = bit
		}
		blocks = append(blocks, block)
		blockNum++
		if blockNum >= Q+K {
			break
		}
	}

	var sum float64 = 0.0
	for blockNumber, eachBlocks := range blocks {
		var _blockNumber_float64 float64 = float64(blockNumber)

		// (2) the L-bit value is used as an index into the table
		var _index_T uint64 = array2BinaryInt(eachBlocks)

		if _blockNumber_float64 < _float64_Q {
			// (2) The block number of the last occurrence of each L-bit block is noted in the table
			T[_index_T] = _blockNumber_float64 + 1.0
		} else {
			// (3) Examine each of the K blocks in the test segment and determine the number of blocks since the last occurrence of the same L-bit block (i.e., i – T[j]).
			sum += math.Log2(float64(blockNumber) + 1.0 - T[_index_T])
			T[_index_T] = float64(blockNumber) + 1.0
		}
	}

	// (4) Compute the test statistic
	var f_n float64 = sum / float64(K)

	// (5) Compute P-value
	var P_value float64 = math.Erfc(math.Abs((f_n - expectedValue_mu[L-1]) / (math.Sqrt2 * variance_sigma[L-1])))

	return P_value, P_value >= 0.01, nil
}
