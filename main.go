package main

import (
	"flag"
	"fmt"
	"os"

	stream "github.com/notJoon/drbg/bitstream"
	nist "github.com/notJoon/drbg/nist"
)

// TODO: pretty print options with detailed output

func main() {
	fmt.Println("NIST Statistical Test Suite")
	allTests := flag.Bool("all", false, "Run all tests")
	frequency := flag.Bool("frequency", false, "Run Frequency (Monobit) Test")
	blockFrequency := flag.Bool("block", false, "Run Frequency Test within a Block")
	runs := flag.Bool("runs", false, "Run Runs Test")
	// longestRun := flag.Bool("longest-run", false, "Run Test for the Longest Run of Ones in a Block")
	// rank := flag.Bool("rank", false, "Run Binary Matrix Rank Test")
	// dft := flag.Bool("dft", false, "Run Discrete Fourier Transform (Spectral) Test")
	// nonOverlappingTemplate := flag.Bool("non-overlapping-template", false, "Run Non-overlapping Template Matching Test")
	// overlappingTemplate := flag.Bool("overlapping-template", false, "Run Overlapping Template Matching Test")
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

	if *allTests || *frequency {
		fmt.Println("Running Frequency (Monobit) Test...")
		p_val, isRandom, err := nist.FrequencyTest(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("p-value: ", p_val)
		fmt.Println("Is random: ", isRandom)
	}
	if *allTests || *blockFrequency {
		fmt.Println("Running Frequency Test within a Block...")
		p_val, isRandom, err := nist.BlockFrequencyTest(bs, 128)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("p-value: ", p_val)
		fmt.Println("Is random: ", isRandom)
	}

	if *allTests || *runs {
		fmt.Println("Running Runs Test...")
		p_val, isRandom, err := nist.Runs(bs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("p-value: ", p_val)
		fmt.Println("Is random: ", isRandom)
	}

	if *verbose {
		fmt.Println("Verbose output enabled")
		// TODO: detailed output
	}
}
