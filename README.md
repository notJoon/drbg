# drgb

Pseudo-Random Number Generator validate testing suite

## Overview

The NIST SP-800-22 test framework provides a standardized method for evaluating the randomness of random and pseudorandom number sequences. Developed by the National Institute of Standards and Technology (NIST), this suite is primarily used in cryptographic applications to assess the quality of random number generators.

### Purpose of the NIST SP-800-22 Tests

These tests statistically assess whether a data sequence has been randomly generated, crucial for generating unpredictable cryptographic keys, initialization vectors, passwords, and other vital security elements.

## How to Use

example REPL command:

```plain
go run main.go -file rand_data/numbers.bin -all -template "0010111011" -block-size 100
```

To use this testing framework, prepare the sequence of data to be tested (The test file should contain at least 1000 data points.), perform each test, and interpret the results to evaluate the adequacy of the random number generator.

Typically, results are labeled **_PASS_** or **_FAIL_** based on their `p-values`; a sequence passes a test if its p-value is greater than `0.01`, indicating decision rules in the document which is the pivot satisfactory randomness.

The suite includes various tests, each examining specific properties or patterns within the data. This includes frequency tests, block frequency tests, runs tests, matrix rank tests, and more, each designed to detect non-random occurrences and ensure the data does not follow predictable patterns.

## List of Tests

The tests include all the tests specified in NIST SP-800-22 document. More detailed explanation of each tests please refer the NIST's document[^1]. The sections and page numbers also refer to this document.

### Frequency (Monobit) Test

> _Section 2.1 p.24_

This test checks the frequency of single bits to ensure that the occurrences of 0s and 1s are approximately the same within the sequence.

### Frequency Test within a Block

> _Section 2.2 p.26_

It evaluates the frequency of bits within a block of data to check for balance, ensuring no uneven distribution of 0s and 1s.

### Runs Test

> _Section 2.3 p.27_

This test measures how frequently runs of consecutive identical bits occur, examining the sequence for high uniformity that might indicate non-randomness.

### Test for the Longest Run of Ones in a Block

> _Section 2.4 p.29_

Determines the length of the longest run of '1's in a specified block, assessing the sequence for unusual patterns.

### Binary Matrix Rank Test

> _Section 2.5 p.32_

Uses the rank of matrices to evaluate the dimensional structure of the data, checking for linear dependencies.

### Discrete Fourier Transform (Spectral) Test

> _Section 2.6 p.34_

Analyzes the sequence in the frequency domain to detect periodic features, which could indicate predictability.

### Non-overlapping Template Matching Test

> _Section 2.7 p.36_

Assesses how frequently certain predefined bit patterns appear within the sequence, checking for their unexpected repetition or rarity.

### Overlapping Template Matching Test

> _Section 2.8 p.39_

Evaluates the frequency of overlapping patterns, looking for deviations from expected randomness.

### Maurer's "Universal Statistical" Test

> _Section 2.9 p.42_

Minimum length of test case: 387,840

Measures the complexity of data sequences to evaluate their randomness, assessing the sequence’s entropy and compression potential.

### Linear Complexity Test

> _Section 2.10 p.46_

Evaluates how complex or simple a sequence can be described, providing insight into the sequence's algorithmic complexity.

### Serial Test

> _Section 2.11, p.48_

Checks for the frequency of repeating patterns within the sequence to detect structured deviations from randomness.

### Approximate Entropy Test

> _Section 2.12 p.51_

Measures the entropy of the sequence, quantifying the degree of randomness or unpredictability.

### Cumulative Sums (Cusum) Test

> _Section 2.13 p.53_

Analyzes whether the cumulative sum of the sequence deviates from expected behavior, indicating potential biases.

### Random Excursions Test

> _Section 2.14 p.55_

Analyzes random deviations from the mean to assess patterns that could suggest non-randomness.

### Random Excursions Variant Test

> _Section 2.15 p.60_

Further examines random excursions using various states, providing additional analysis on deviations from randomness.

## Reference

[^1]: [A Stastical Test Suite for Random and Pseudorandom Number Generators for Cryptographic Applications](<https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-22r1a.pdf>)
