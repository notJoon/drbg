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
	blockFrequency := flag.Bool("block", false, "Run Frequency Test within a Block")
	runs := flag.Bool("runs", false, "Run Runs Test")
	longestRun := flag.Bool("longest-run", false, "Run Test for the Longest Run of Ones in a Block")
	rank := flag.Bool("rank", false, "Run Binary Matrix Rank Test")
	dft := flag.Bool("dft", false, "Run Discrete Fourier Transform (Spectral) Test")

	nonOverlappingTemplate := flag.Bool("non-overlapping", false, "Run Non-overlapping Template Matching Test")
	overlappingTemplate := flag.Bool("overlapping", false, "Run Overlapping Template Matching Test")
	// specifies the template B to match. Must be string of ones and zeros (e.g. "001")
	templateB := flag.String("template", "", "The template B to be matched (a string of ones and zeros)")
	// specified the length of the substrting to test, in bits.
	blockSize := flag.Uint64("block-size", 0, "The length in bits of the substring to be tested")

	// universal := flag.Bool("universal", false, "Run Maurer's Universal Statistical Test")
	// linearComplexity := flag.Bool("linear-complexity", false, "Run Linear Complexity Test")
	// serial := flag.Bool("serial", false, "Run Serial Test")
	// approximateEntropy := flag.Bool("approximate-entropy", false, "Run Approximate Entropy Test")
	// cusum := flag.Bool("cusum", false, "Run Cumulative Sums (Cusums) Test")
	// randomExcursions := flag.Bool("random-excursions", false, "Run Random Excursions Test")
	// randomExcursionsVariant := flag.Bool("random-excursions-variant", false, "Run Random Excursions Variant Test")

	filename := flag.String("file", "", "File containing the random bits")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
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

	bs, err := stream.FromFile(*filename)
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
		p_val, isRandom, err := nist.FrequencyTest(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if isRandom {
			pass++
			t.AppendRow([]interface{}{"Frequency (Monobit) Test", p_val, "Pass"})
		} else {
			fail++
			t.AppendRow([]interface{}{"Frequency (Monobit) Test", p_val, "Fail"})
		}
	}
	if *allTests || *blockFrequency {
		p_val, isRandom, err := nist.BlockFrequencyTest(bs, 128)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if isRandom {
			pass++
			t.AppendRow([]interface{}{"Frequency Test within a Block", p_val, "Pass"})
		} else {
			fail++
			t.AppendRow([]interface{}{"Frequency Test within a Block", p_val, "Fail"})
		}
	}

	if *allTests || *runs {
		p_val, isRandom, err := nist.Runs(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if isRandom {
			pass++
			t.AppendRow([]interface{}{"Runs Test", p_val, "Pass"})
		} else {
			fail++
			t.AppendRow([]interface{}{"Runs Test", p_val, "Fail"})
		}
	}

	if *allTests || *longestRun {
		p_val, isRandom, err := nist.LongestRunOfOnes(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if isRandom {
			pass++
			t.AppendRow([]interface{}{"Test for the Longest Run of Ones in a Block", p_val, "Pass"})
		} else {
			fail++
			t.AppendRow([]interface{}{"Test for the Longest Run of Ones in a Block", p_val, "Fail"})
		}
	}

	if *allTests || *rank {
		p_val, isRandom, err := nist.Rank(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if isRandom {
			pass++
			t.AppendRow([]interface{}{"Binary Matrix Rank Test", p_val, "Pass"})
		} else {
			fail++
			t.AppendRow([]interface{}{"Binary Matrix Rank Test", p_val, "Fail"})
		}
	}

	if *allTests || *dft {
		p_val, isRandom, err := nist.DFT(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if isRandom {
			pass++
			t.AppendRow([]interface{}{"Discrete Fourier Transform (Spectral) Test", p_val, "Pass"})
		} else {
			fail++
			t.AppendRow([]interface{}{"Discrete Fourier Transform (Spectral) Test", p_val, "Fail"})
		}
	}

	if *allTests || *nonOverlappingTemplate {
		if *templateB == "" {
			fmt.Println("Error: template B is required for Non-overlapping Template Matching Test.\nUse -template \"001\" (or other tmeplate)")
			os.Exit(1)
		}
		if *blockSize == 0 {
			fmt.Println("Error: block size is required for Non-overlapping Template Matching Test.\nUse -block-size 10 (or other block size)")
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
				fmt.Printf("Error: invalid character in template B: %c\n", c)
				os.Exit(1)
			}
		}
		p_value, isRandom, err := nist.NonOverlappingTemplateMatching(B, *blockSize, bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if isRandom {
			pass++
			t.AppendRow([]interface{}{"Non-overlapping Template Matching Test", p_value, "Pass"})
		} else {
			fail++
			t.AppendRow([]interface{}{"Non-overlapping Template Matching Test", p_value, "Fail"})
		}
	}

	if *allTests || *overlappingTemplate {
		if *templateB == "" {
			fmt.Println("Error: template B is required for Non-overlapping Template Matching Test.\nUse -template \"001\" (or other tmeplate)")
			os.Exit(1)
		}
		if *blockSize == 0 {
			fmt.Println("Error: block size is required for Non-overlapping Template Matching Test.\nUse -block-size 10 (or other block size)")
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
				fmt.Printf("Error: invalid character in template B: %c\n", c)
				os.Exit(1)
			}
		}
		p_value, isRandom, err := nist.OverlappingTemplateMatching(B, *blockSize, bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if isRandom {
			pass++
			t.AppendRow([]interface{}{"Overlapping Template Matching Test", p_value, "Pass"})
		} else {
			fail++
			t.AppendRow([]interface{}{"Overlapping Template Matching Test", p_value, "Fail"})
		}
	}

	if *verbose {
		fmt.Println("Verbose output enabled")
		// TODO: detailed output
	}

	t.AppendFooter(table.Row{"", "Total Tests", pass + fail})
	t.AppendFooter(table.Row{"", "Pass", pass})
	t.AppendFooter(table.Row{"", "Fail", fail})
	t.Render()
}
