package bitstream

import (
	"errors"
	"reflect"
	"testing"
)

func TestBitStreamReader(t *testing.T) {
	tests := []struct {
		name        string
		initData    []byte
		readOps     []func(*BitStreamReader) ([]byte, error)
		wantResults [][]byte
		wantErr     []error
	}{
		{
			name:     "Read bits successfully",
			initData: []byte{0b10101111, 0b11000001}, // 10101111 11000001
			readOps: []func(*BitStreamReader) ([]byte, error){
				func(r *BitStreamReader) ([]byte, error) { return r.Read(10) }, // Read first 10 bits
			},
			wantResults: [][]byte{
				{0b10101111, 0b11000000}, // 10101111 11000000
			},
			wantErr: []error{nil},
		},
		{
			name:     "Read single byte",
			initData: []byte{0xAA, 0x55},
			readOps: []func(*BitStreamReader) ([]byte, error){
				func(r *BitStreamReader) ([]byte, error) {
					bs, err := r.ReadByte()
					if err != nil {
						return nil, err
					}
					return []byte{bs}, nil
				},
			},
			wantResults: [][]byte{
				{0xAA},
			},
			wantErr: []error{nil},
		},
		{
			name:     "Read bits with not enough bits left",
			initData: []byte{0xFF}, // 11111111
			readOps: []func(*BitStreamReader) ([]byte, error){
				func(r *BitStreamReader) ([]byte, error) { return r.Read(9) }, // Attempt to read 9 bits from 8 bits available
			},
			wantResults: [][]byte{
				nil,
			},
			wantErr: []error{ErrNotEnoughBits},
		},
		{
			name:     "Read bit by bit",
			initData: []byte{0b01010101}, // 01010101
			readOps: []func(*BitStreamReader) ([]byte, error){
				func(r *BitStreamReader) ([]byte, error) {
					b, err := r.ReadBit()
					if err != nil {
						return nil, err
					}
					return []byte{b}, nil
				},
				func(r *BitStreamReader) ([]byte, error) {
					b, err := r.ReadBit()
					if err != nil {
						return nil, err
					}
					return []byte{b}, nil
				},
				func(r *BitStreamReader) ([]byte, error) {
					b, err := r.ReadBit()
					if err != nil {
						return nil, err
					}
					return []byte{b}, nil
				},
				func(r *BitStreamReader) ([]byte, error) {
					b, err := r.ReadBit()
					if err != nil {
						return nil, err
					}
					return []byte{b}, nil
				},
				func(r *BitStreamReader) ([]byte, error) {
					b, err := r.ReadBit()
					if err != nil {
						return nil, err
					}
					return []byte{b}, nil
				},
			},
			wantResults: [][]byte{
				{0}, {1}, {0}, {1}, {0},
			},
			wantErr: []error{nil, nil, nil, nil, nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitstream(tt.initData)
			reader := NewBitStreamReader(bs)
			for i, op := range tt.readOps {
				result, err := op(reader)
				if !errors.Is(err, tt.wantErr[i]) {
					t.Errorf("Test '%s': Expected error %v, got %v", tt.name, tt.wantErr[i], err)
				}
				if err == nil && !reflect.DeepEqual(result, tt.wantResults[i]) {
					t.Errorf("Test '%s': Expected result %v, got %v", tt.name, tt.wantResults[i], result)
				}
			}
		})
	}
}
