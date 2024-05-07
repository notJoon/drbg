package main

import (
	"flag"
	"fmt"
	"os"

	stream "github.com/notJoon/drbg/bitstream"
	nist "github.com/notJoon/drbg/nist"

	"github.com/jedib0t/go-pretty/table"
)

func main() {
	allTests := flag.Bool("all", false, "Run all tests")

	frequency := flag.Bool("frequency", false, "Run Frequency (Monobit) Test")
	blockFrequency := flag.Bool("block-frequency", false, "Run Frequency Test within a Block")
	blockFrequencyBlockSize := flag.Uint64("frequency-block-size", 128, "The length in bits of the substring to be tested")

	runs := flag.Bool("runs", false, "Run Runs Test")
	longestRun := flag.Bool("longest-run", false, "Run Test for the Longest Run of Ones in a Block")

	rank := flag.Bool("rank", false, "Run Binary Matrix Rank Test")
	dft := flag.Bool("dft", false, "Run Discrete Fourier Transform (Spectral) Test")

	nonOverlappingTemplate := flag.Bool("non-overlapping", false, "Run Non-overlapping Template Matching Test.\nDefault template is \"000000001\" and block size is 10 bits.")
	overlappingTemplate := flag.Bool("overlapping", false, "Run Overlapping Template Matching Test.\nDefault template is \"000000001\" and block size is 10 bits.")
	// specifies the template B to match. Must be string of ones and zeros (e.g. "001")
	templateB := flag.String("template", "000000001", "The template B to be matched (a string of ones and zeros)")
	// specified the length of the substrting to test, in bits.
	blockSize := flag.Uint64("block-size", 20, "The length in bits of the substring to be tested")

	universal := flag.Bool("universal", false, "Run Maurer's Universal Statistical Test")

	linearComplexity := flag.Bool("linear", false, "Run Linear Complexity Test. Default block size is 500 bits.")
	inputSize := flag.Uint64("m", 500, "The length of the block to be tested")

	serial := flag.Bool("serial", false, "Run Serial Test. Default block size is 16 bits.")
	serialBlockSize := flag.Uint64("serial-size", 16, "The length in bits of the substring to be tested")

	approximateEntropy := flag.Bool("entropy", false, "Run Approximate Entropy Test")
	approximateEntropyBlockSize := flag.Uint64("entropy-block-size", 10, "The length in bits of the substring to be tested")

	cusum := flag.Bool("cusum", false, "Run Cumulative Sums (Cusums) Test. Default mode is 0 (forward).")
	mode := flag.Int("mode", 0, "The mode of the test (0 or 1)")

	randomExcursions := flag.Bool("random-excursions", false, "Run Random Excursions Test")
	randomExcursionsVariant := flag.Bool("random-excursions-variant", false, "Run Random Excursions Variant Test")

	filename := flag.String("file", "", "File containing the random bits")

	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *filename == "" {
		fmt.Println("Error: No file specified")
		os.Exit(1)
	}

	var (
		bs  *stream.BitStream
		err error
	)

	// regulation of the bitstream
	// ????
	if *frequency {
		bs, err = stream.FromFileWithLimit(*filename, 100)
	} else {
		bs, err = stream.FromFile(*filename)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// test result counters
	pass, fail := 0, 0

	// Draw table for test results
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"NIST Statistical Test Suite", "p-value", "Result"})

	if *allTests || *frequency {
		testName := "Frequency (Monobit) Test"
		p_val, isRandom, err := nist.FrequencyTest(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		writeResult(t, testName, p_val, isRandom, &pass, &fail)
	}
	if *allTests || *blockFrequency {
		testName := "Frequency Test within a Block"
		p_val, isRandom, err := nist.BlockFrequencyTest(bs, uint64(*blockFrequencyBlockSize))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		writeResult(t, testName, p_val, isRandom, &pass, &fail)
	}

	if *allTests || *runs {
		testName := "Runs Test"
		p_val, isRandom, err := nist.Runs(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		writeResult(t, testName, p_val, isRandom, &pass, &fail)
	}

	if *allTests || *longestRun {
		testName := "Test for the Longest Run of Ones in a Block"
		p_val, isRandom, err := nist.LongestRunOfOnes(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		writeResult(t, testName, p_val, isRandom, &pass, &fail)
	}

	if *allTests || *rank {
		testName := "Binary Matrix Rank Test"
		p_val, isRandom, err := nist.Rank(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		writeResult(t, testName, p_val, isRandom, &pass, &fail)
	}

	if *allTests || *dft {
		testName := "Discrete Fourier Transform (Spectral) Test"
		p_val, isRandom, err := nist.DFT(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		writeResult(t, testName, p_val, isRandom, &pass, &fail)
	}

	if *allTests || *nonOverlappingTemplate {
		testName := "Non-overlapping Template Matching Test"
		if *templateB == "" {
			fmt.Println("Error (non-overlapping template test): template B is required for Non-overlapping Template Matching Test.\nUse -template \"001\" (or other tmeplate)")
			os.Exit(1)
		}
		if *blockSize == 0 {
			fmt.Println("Error (non-overlapping template test): block size is required for Non-overlapping Template Matching Test.\nUse -block-size 10 (or other block size)")
			os.Exit(1)
		}
		B := make([]uint8, len(*templateB))
		for i, c := range *templateB {
			switch c {
			case '0':
				B[i] = 0
			case '1':
				B[i] = 1
			default:
				fmt.Printf("Error (non-overlapping template test): invalid character in template B: %c\n", c)
				os.Exit(1)
			}
		}
		p_value, isRandom, err := nist.NonOverlappingTemplateMatching(B, *blockSize, bs)
		if err != nil {
			fmt.Printf("Error (non-overlapping template test): %v\n", err)
			os.Exit(1)
		}

		writeResult(t, testName, p_value, isRandom, &pass, &fail)
	}

	if *allTests || *overlappingTemplate {
		testName := "Overlapping Template Matching Test"
		if *templateB == "" {
			fmt.Println("Error (overlapping templelate test): template B is required for Non-overlapping Template Matching Test.\nUse -template \"001\" (or other tmeplate)")
			os.Exit(1)
		}
		if *blockSize == 0 {
			fmt.Println("Error (overlapping templelate test): block size is required for Non-overlapping Template Matching Test.\nUse -block-size 10 (or other block size)")
			os.Exit(1)
		}
		B := make([]uint8, len(*templateB))
		for i, c := range *templateB {
			switch c {
			case '0':
				B[i] = 0
			case '1':
				B[i] = 1
			default:
				fmt.Printf("Error (overlapping templelate test): invalid character in template B: %c\n", c)
				os.Exit(1)
			}
		}
		p_value, isRandom, err := nist.OverlappingTemplateMatching(B, *blockSize, bs)
		if err != nil {
			fmt.Printf("Error (overlapping templelate test): %v\n", err)
			os.Exit(1)
		}

		writeResult(t, testName, p_value, isRandom, &pass, &fail)
	}

	if *universal {
		testName := "Maurer's Universal Statistical Test"
		p_val, isRandom, err := nist.UniversalRecommendedValues(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		writeResult(t, testName, p_val, isRandom, &pass, &fail)
	}

	if *allTests || *linearComplexity {
		testName := "Linear Complexity Test"
		if *inputSize < 500 || *inputSize > 5000 {
			fmt.Println("Error: input size must be between 500 and 5000")
			os.Exit(1)
		}
		p_val, isRandom, err := nist.LinearComplexity(*inputSize, bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		writeResult(t, testName, p_val, isRandom, &pass, &fail)
	}

	if *allTests || *serial {
		testName := "Serial Test"
		p_val, isRandom, err := nist.Serial(*serialBlockSize, bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if isRandom != nil {
			pass++
			t.AppendRow([]interface{}{testName, fmt.Sprintf("%.2f", p_val), "Pass"})
		} else {
			fail++
			t.AppendRow([]interface{}{testName, fmt.Sprintf("%.2f", p_val), "Fail"})
		}
	}

	if *allTests || *approximateEntropy {
		testName := "Approximate Entropy Test"
		p_val, isRandom, err := nist.ApproximateEntropy(*approximateEntropyBlockSize, bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		writeResult(t, testName, p_val, isRandom, &pass, &fail)
	}

	if *allTests || *cusum {
		testName := "Cumulative Sums Test"
		p_val, isRandom, err := nist.CumulativeSums(*mode, bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		writeResult(t, testName, p_val, isRandom, &pass, &fail)
	}

	if *allTests || *randomExcursions {
		testName := "Random Excursions Test"
		p_val, isRandom, err := nist.RandomExcursions(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if isRandom != nil {
			pass++
			t.AppendRow([]interface{}{testName, fmt.Sprintf("%.2f", p_val), "Pass"})
		} else {
			fail++
			t.AppendRow([]interface{}{testName, fmt.Sprintf("%.2f", p_val), "Fail"})
		}
	}

	if *allTests || *randomExcursionsVariant {
		testName := "Random Excursions Variant Test"
		p_val, isRandom, err := nist.RandomExcursionsVariant(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if isRandom != nil {
			pass++
			t.AppendRow([]interface{}{testName, fmt.Sprintf("%.2f", p_val), "Pass"})
		} else {
			fail++
			t.AppendRow([]interface{}{testName, fmt.Sprintf("%.2f", p_val), "Fail"})
		}
	}

	t.AppendFooter(table.Row{"", "Total Tests", pass + fail})
	t.AppendFooter(table.Row{"", "Pass", pass})
	t.AppendFooter(table.Row{"", "Fail", fail})
	t.Render()
}

// writeResult writes the result of a test to the table
func writeResult(t table.Writer, testName string, pValue float64, isRandom bool, pass *int, fail *int) {
	result := "Fail"
	if isRandom {
		result = "Pass"
		*pass += 1
	} else {
		*fail += 1
	}
	t.AppendRow([]interface{}{testName, fmt.Sprintf("%.2f", pValue), result})
}
