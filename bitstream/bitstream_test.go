package bitstream

import (
	"testing"
)

func TestBitStream(t *testing.T) {
	tests := []struct {
		name     string
		initData []byte
		length   int
		ops      []func(*BitStream) error
		wantBits []byte
		wantErr  error
	}{
		{
			name:     "Empty BitStream",
			initData: []byte{},
			length:   0,
			ops:      nil,
			wantBits: []byte{},
			wantErr:  nil,
		},
		{
			name:     "Single Byte No Change",
			initData: []byte{0xAA},
			length:   8,
			ops:      nil,
			wantBits: []byte{0xAA},
			wantErr:  nil,
		},
		{
			name:     "Single Bit Set",
			initData: []byte{0x00},
			length:   8,
			ops: []func(*BitStream) error{
				func(bs *BitStream) error { return bs.SetBit(0, 1) },
			},
			wantBits: []byte{0x80},
			wantErr:  nil,
		},
		{
			name:     "Read Bit",
			initData: []byte{0x80},
			length:   8,
			ops: []func(*BitStream) error{
				func(bs *BitStream) error {
					bit, err := bs.Bit(0)
					if err != nil || bit != 1 {
						t.Fatalf("Expected bit 1, got %d, err: %v", bit, err)
					}
					return nil
				},
			},
			wantBits: []byte{0x80},
			wantErr:  nil,
		},
		{
			name:     "Set and Read Bit",
			initData: []byte{0x00},
			length:   8,
			ops: []func(*BitStream) error{
				func(bs *BitStream) error { return bs.SetBit(7, 1) },
				func(bs *BitStream) error {
					bit, err := bs.Bit(7)
					if err != nil || bit != 1 {
						t.Fatalf("Expected bit 1, got %d, err: %v", bit, err)
					}
					return nil
				},
			},
			wantBits: []byte{0x01},
			wantErr:  nil,
		},
		{
			name:     "Append Bit",
			initData: []byte{0x80}, // 10000000
			length:   8,
			ops: []func(*BitStream) error{
				func(bs *BitStream) error { return bs.Append(1) }, // Append 1
			},
			wantBits: []byte{0x80, 0x80}, // 10000000 10000000
			wantErr:  nil,
		},
		{
			name:     "Invalid Index",
			initData: []byte{0x00},
			length:   8,
			ops: []func(*BitStream) error{
				func(bs *BitStream) error { _, err := bs.Bit(8); return err },
			},
			wantBits: []byte{0x00},
			wantErr:  ErrOutOfRange,
		},
		{
			name:     "Invalid Bit Value",
			initData: []byte{0x00},
			length:   8,
			ops: []func(*BitStream) error{
				func(bs *BitStream) error { return bs.Append(2) }, // Invalid bit value
			},
			wantBits: []byte{0x00},
			wantErr:  ErrInvalidBitValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitStream(tt.initData)
			var err error
			for _, op := range tt.ops {
				if err = op(bs); err != nil {
					break
				}
			}

			if tt.wantErr != nil && err != tt.wantErr {
				t.Errorf("Expected error %v, got %v", tt.wantErr, err)
			} else if tt.wantErr == nil && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(bs.Bytes()) != len(tt.wantBits) {
				t.Fatalf("Expected bytes %v, got %v", tt.wantBits, bs.Bytes())
			}

			for i, b := range bs.Bytes() {
				if b != tt.wantBits[i] {
					t.Errorf("Expected byte at index %d to be %b, got %b", i, tt.wantBits[i], b)
				}
			}
		})
	}
}

func TestBit(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		index    int
		expected byte
		err      error
	}{
		{
			name:     "All zeros",
			data:     []byte{0x00, 0x00}, // 00000000 00000000
			index:    0,
			expected: 0,
			err:      nil,
		},
		{
			name:     "All ones",
			data:     []byte{0xff, 0xff}, // 11111111 11111111
			index:    7,
			expected: 1,
			err:      nil,
		},
		{
			name:     "Mixed pattern",
			data:     []byte{0xaa, 0x55}, // 10101010 01010101
			index:    8,
			expected: 0,
			err:      nil,
		},
		{
			name:     "Mixed pattern",
			data:     []byte{0xaa, 0x55}, // 10101010 01010101
			index:    9,
			expected: 1,
			err:      nil,
		},
		{
			name:     "Index out of range - negative",
			data:     []byte{0x00},
			index:    -1,
			expected: 0,
			err:      ErrOutOfRange,
		},
		{
			name:     "Index out of range - too large",
			data:     []byte{0x00, 0x00},
			index:    16,
			expected: 0,
			err:      ErrOutOfRange,
		},
	}

	for _, tt := range tests {
		bs := NewBitStream(tt.data)
		result, err := bs.Bit(tt.index)
		if err != tt.err {
			t.Errorf("%s: expected error %v, got %v", tt.name, tt.err, err)
		}
		if result != tt.expected {
			t.Errorf("%s: expected result %d, got %d", tt.name, tt.expected, result)
		}
	}
}
