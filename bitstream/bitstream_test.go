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
