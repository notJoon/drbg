package nist

import (
	"math"
	"testing"

	b "github.com/notJoon/drbg/bitstream"
)

func TestMonobit_FrequenctTest(t *testing.T) {
	epsilon := 0.0001
	tests := []struct {
		name             string
		bs               *b.BitStream
		expectedP        float64
		expectedIsRandom bool
		expectedErr      bool
	}{
		{
			name:             "Empty bitstream",
			bs:               b.NewBitStream([]byte{}),
			expectedP:        0,
			expectedIsRandom: false,
			expectedErr:      true,
		},
		{
			name:             "All zeros",
			bs:               b.NewBitStream([]byte{0x00, 0x00, 0x00, 0x00}),
			expectedP:        0, // p-value should be exteremely close to 0
			expectedIsRandom: false,
			expectedErr:      false,
		},
		{
			name:             "All ones",
			bs:               b.NewBitStream([]byte{0xFF, 0xFF, 0xFF, 0xFF}),
			expectedP:        0,
			expectedIsRandom: false,
			expectedErr:      false,
		},
		{
			name:             "Equal ones and zeros",
			bs:               b.NewBitStream([]byte{0xAA, 0xAA, 0xAA, 0xAA}),
			expectedP:        1,
			expectedIsRandom: true,
			expectedErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, isRandom, err := FrequencyTest(tt.bs)
			if (err != nil) != tt.expectedErr {
				t.Errorf("FrequencyTest() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
			if !tt.expectedErr && (!almostEq(p, tt.expectedP, epsilon) || isRandom != tt.expectedIsRandom) {
				t.Errorf("FrequencyTest() = %v, %v, expected %v, %v", p, isRandom, tt.expectedP, tt.expectedIsRandom)
			}
		})
	}
}

// almostEq checks if two floating-point numbers are close enough.
func almostEq(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}
