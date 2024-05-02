package bitstream

import "errors"

var ErrNotEnoughBits = errors.New("not enough bits")

// BitStreamReader provides a way to read bits from a Bitstream.
type BitStreamReader struct {
	bs     *BitStream // Bitstream to read from
	offset int        // current offset in the Bitstream
}

// NewBitStreamReader creates a new BitStreamReader from the provided Bitstream.
func NewBitStreamReader(bs *BitStream) *BitStreamReader {
	return &BitStreamReader{bs: bs}
}

// Read reads the specified number of bits from the Bitstream
// and returns them as a byte slice.
//
// It returns an error if there are not enough bits to read.
func (r *BitStreamReader) Read(n int) ([]byte, error) {
	if r.offset+n >= r.bs.Len() {
		return nil, ErrNotEnoughBits
	}

	data := make([]byte, (n+msbIndex)/bitSize)
	bitcount := 0
	for i := 0; i < n; i++ {
		bit, _ := r.bs.Bit(r.offset + i)
		data[bitcount/bitSize] |= bit << uint(msbIndex-bitcount%bitSize)
		bitcount++
	}
	r.offset += n
	return data, nil
}

// ReadByte reads a 8 bits from the Bitstream and returns them as a byte.
// It returns an error if there are not enough bits to read.
func (r *BitStreamReader) ReadByte() (byte, error) {
	bs, err := r.Read(8)
	if err != nil {
		return 0, err
	}
	return bs[0], nil
}

// ReadBit reads a single bit from the BitStream and returns it as a byte.
// It returns an error if there are not enough bits to read.
func (r *BitStreamReader) ReadBit() (byte, error) {
	b, err := r.bs.Bit(r.offset)
	if err != nil {
		return 0, err
	}
	r.offset++
	return b, nil
}
