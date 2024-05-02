package bitstream

import (
	"errors"
	"reflect"
	"testing"
)

func TestBitStreamWriter(t *testing.T) {
	tests := []struct {
		name          string
		initData      []byte
		writeOps      []func(*BitStreamWriter) error
		flushExpected []byte
		wantErr       []error
	}{
		{
			name:     "Write bytes to stream",
			initData: []byte{},
			writeOps: []func(*BitStreamWriter) error{
				func(w *BitStreamWriter) error { return w.Write([]byte{0xAA, 0x55}) },
			},
			flushExpected: []byte{0xAA, 0x55},
			wantErr:       []error{nil},
		},
		{
			name:     "Write single bit",
			initData: []byte{},
			writeOps: []func(*BitStreamWriter) error{
				func(w *BitStreamWriter) error { return w.WriteBit(1) },
				func(w *BitStreamWriter) error { return w.WriteBit(0) },
			},
			flushExpected: []byte{0x80}, // 10000000
			wantErr:       []error{nil, nil},
		},
		{
			name:     "Invalid bit value",
			initData: []byte{},
			writeOps: []func(*BitStreamWriter) error{
				func(w *BitStreamWriter) error { return w.WriteBit(2) }, // Invalid bit value
			},
			flushExpected: nil,
			wantErr:       []error{ErrInvalidBitValue},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitStream(tt.initData)
			writer := NewBitStreamWriter(bs)
			var err error
			for i, op := range tt.writeOps {
				if err = op(writer); err != nil {
					if !errors.Is(err, tt.wantErr[i]) {
						t.Errorf("Test '%s': Expected error %v, got %v", tt.name, tt.wantErr[i], err)
					}
					return
				}
			}
			writer.Flush()
			if !reflect.DeepEqual(bs.Bytes(), tt.flushExpected) {
				t.Errorf("Test '%s': Expected bytes %v after flush, got %v", tt.name, tt.flushExpected, bs.Bytes())
			}
		})
	}
}
