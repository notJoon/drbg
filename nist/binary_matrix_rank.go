package nist

import (
	"math"

	b "github.com/notJoon/drbg/bitstream"
)

// ref: A Statistical Test Suite for Random and Pseudorandom Number Generators for Cryptographic Application
// Section 2.5 Binary Matrix Rank Test (p. 32)

// Rank performs the Binary Matrix Rank Test on the given bitstream.
// The purpose of this test is to check for linear dependence among fixed-length substrings of the original sequence.
// The test constructs matrices from the input sequence and determines the rank of each matrix.
// It then compares the rank distribution to the expected distribution for a random sequence.
// Deviations from the expected distribution indicate non-randomness.
//
// The test proceeds as follows:
//  1. Divide the input sequence into M = n/Q non-overlapping blocks of length Q, where n is the length of the input sequence.
//  2. Construct a Q x Q matrix from each block by writing the bits in the block in row-major order.
//  3. Determine the binary rank of each matrix using Gaussian elimination.
//  4. Count the number of matrices with each rank (full rank, full rank-1, and other ranks).
//  5. Compute the Chi-square statistic based on the observed and expected counts for each rank.
//  6. Compute the p-value using the incomplete Gamma function.
func Rank(bs *b.BitStream) (float64, bool, error) {
	var (
		M uint64 = 32 // number of rows in the matrix
		Q uint64 = 32 // number of columns in the matrix
	)

	// sequentially divide the sequence into M*Q bit disjoint blocks
	n := bs.Len()
	N := uint64(n) / (M * Q)
	R := make([]uint64, N)   // TODO: change name to rank
	F := make([]uint64, M+1) // number of matrices with rank_i = index (index means, rank)

	eposilonIdx := 0
	matrices := make([][][]uint8, N)

	for i := 0; i < int(N); i++ {
		matrices[i] = make([][]uint8, M)
		for j := 0; j < int(M); j++ {
			matrices[i][j] = make([]uint8, Q)
			for k := 0; k < int(Q); k++ {
				matrices[i][j][k], _ = bs.Bit(eposilonIdx)
				eposilonIdx++
			}
		}
	}

	for i, matrix := range matrices {
		// determine the binary rank of each matrix (where l := 1, ..., N)
		// ref: Appendix A. (page 33)
		R[i] = RankComputationOfBinaryMatrices(matrix)

		// let F_M = number of matrices with R_l = M (full rank)
		F[R[i]]++
	}

	// compute chi-square value
	var (
		__F_M_float64         = float64(F[M])
		__F_M_minus_1_float64 = float64(F[M-1])
		__N_float64           = float64(N)
		chi_square            = (__F_M_float64-0.2888*__N_float64)*(__F_M_float64-0.2888*__N_float64)/(0.2888*__N_float64) + (__F_M_minus_1_float64-0.5776*__N_float64)*(__F_M_minus_1_float64-0.5776*__N_float64)/(0.5776*__N_float64) + (__N_float64-__F_M_float64-__F_M_minus_1_float64-0.1336*__N_float64)*(__N_float64-__F_M_float64-__F_M_minus_1_float64-0.1336*__N_float64)/(0.1336*__N_float64)
	)

	// compute P-value
	P_value := math.Pow(math.E, -1*chi_square/2)

	return P_value, P_value >= 0.01, nil
}

func RankComputationOfBinaryMatrices(matrix [][]uint8) uint64 {
	// Forward Application of Elementary Row Operations
	// Declare Variables
	var row, col int
	var m int = len(matrix)

	// Step 1. Set i = 1
	i := 0

	// Step 2. If element a(i,i) = 0 (i.e., the element on the diagonal ≠ 1),
	// then swap all elements in the ith row with all elements in the next row that contains a one in the i-th column.
	// (i.e., this row is the kth row, where i < k <= m)
	// If no row contains a “1” in this position, go to step 4.
Forward_STEP2:
	if matrix[i][i] == 0 {
		var tempIndex int
		var isContained bool = false
		for tempIndex = i; tempIndex < m; tempIndex++ {
			if matrix[tempIndex][i] == 1 {
				matrix[i], matrix[tempIndex] = matrix[tempIndex], matrix[i]
				isContained = true
				break
			}
		}
		if !isContained {
			goto Forward_STEP4
		}
	}

	// Step 3. If element a(i,i) = 1, then if any subsequent row contains a “1” in the i-th column,
	// replace each element in that row with the exclusive-OR of that element and the corresponding element in the i-th row.

	// Step 3-a.
	row = i + 1
	// Step 3-b.
Forward_STEP_3B:
	col = i
	// Step 3-c.
	if matrix[row][col] == 0 {
		goto Forward_STEP_3G
	}
	// Step 3-d.
Forward_STEP_3D:
	matrix[row][col] = matrix[row][col] ^ matrix[i][col]
	// Step 3-e.
	if col == (m - 1) {
		goto Forward_STEP_3G
	}
	// Step 3-f.
	col = col + 1
	goto Forward_STEP_3D
	// Step 3-g.
Forward_STEP_3G:
	if row == (m - 1) {
		goto Forward_STEP4
	}
	// Step 3-h.
	row = row + 1
	goto Forward_STEP_3B

	// Step 4.
Forward_STEP4:
	if i < m-2 {
		i = i + 1
		goto Forward_STEP2
	}

	// Step 5. Forward row operations completed.

	// The Subsequent Backward Row Operations
	// Step 1. Set i = m	// But, matrix index [0, m-1].
	i = m - 1

	// Step 2. If element a(i, i) = 0,
Backward_STEP_2:
	if matrix[i][i] == 0 {
		// swap all elements in the i-th row with all elements in the next row that contains a one in the i-th column
		var tempIndex int
		var isContained bool = false
		for tempIndex = i; tempIndex >= 0; tempIndex-- {
			if matrix[tempIndex][i] == 1 {
				matrix[i], matrix[tempIndex] = matrix[tempIndex], matrix[i]
				isContained = true
				break
			}
		}
		if !isContained {
			goto Backward_STEP_4
		}
	}

	// If element a(i, i) = 1,
	// Step 3-a
	row = i - 1

	// Step 3-b
Backward_STEP_3B:
	col = i

	// Step 3-c
	if matrix[row][col] == 0 {
		goto Backward_STEP_3G
	}

	// Step 3-d
Backward_STEP_3D:
	matrix[row][col] = matrix[row][col] ^ matrix[i][col]

	// Step 3-e
	if col == 1 {
		goto Backward_STEP_3G
	}

	// Step 3-f
	col = col - 1
	goto Backward_STEP_3D

	// Step 3-g
Backward_STEP_3G:
	if row == 1 {
		goto Backward_STEP_4
	}

	// Step 3-h.
	row = row - 1
	goto Backward_STEP_3B

	// Step 4.
Backward_STEP_4:
	if i > 2 {
		i = i - 1
		goto Backward_STEP_2
	}

	// Step 5. Backward row operation complete.

	// The rank of the matrix = the number of non-zero rows.
	var rank uint64 = 0
	for _, row := range matrix {
		for _, eachValue := range row {
			if eachValue == 1 {
				rank++
				break
			}
		}
	}

	return rank
}
