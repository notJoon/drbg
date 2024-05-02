package bitstream

import (
	"bufio"
	"errors"
	"os"
	"strconv"
)

const (
	bitSize  = 8
	msbIndex = bitSize - 1
)

var (
	ErrOutOfRange      = errors.New("index out of range")
	ErrInvalidBitValue = errors.New("invalid bit value")
)

// BitStream represents a sequence of bits storedd in a byte slice.
type BitStream struct {
	data []byte // byte slice to store bits
	len  int    // number of bits in the bitstream (length of the byte slice * 8)
}

// NewBitStream creates a new Bitstream from the provided byte slice.
func NewBitStream(data []byte) *BitStream {
	return &BitStream{data: data, len: len(data) * bitSize}
}

// Len returns the length of the bitsream in bits.
func (bs *BitStream) Len() int {
	return bs.len
}

// Bit returns the bit value at the sepcified index.
// It returns an error if the index is out of range.
func (bs *BitStream) Bit(index int) (byte, error) {
	if index < 0 || index >= bs.len {
		return 0, ErrOutOfRange
	}

	byteIndex, bitIndex := getIndexes(index)
	return (bs.data[byteIndex] >> uint(msbIndex-bitIndex)) & 1, nil
}

// SetBit sets the bit value at the specified index.
// It returns an error if the index is out of range.
func (bs *BitStream) SetBit(index int, bit byte) error {
	if index < 0 || index >= bs.len {
		return ErrOutOfRange
	}

	byteIndex, bitIndex := getIndexes(index)
	mask := byte(1 << uint(msbIndex-bitIndex))
	if bit == 1 {
		bs.data[byteIndex] |= mask
	} else {
		bs.data[byteIndex] &= ^mask
	}
	return nil
}

// Append appends a bit to the end of the bitstream.
// It returns an error if the provided bit value is not 0 or 1.
func (bs *BitStream) Append(bit byte) error {
	if bit != 0 && bit != 1 {
		return ErrInvalidBitValue
	}

	byteIndex, bitIndex := getIndexes(bs.len)
	if bitIndex == 0 {
		bs.data = append(bs.data, 0)
	}

	if bit == 1 {
		bs.data[byteIndex] |= 1 << uint(msbIndex-bitIndex)
	}

	bs.len++
	return nil
}

// Bytes returns the underlying byte slice of the bitstream.
func (bs *BitStream) Bytes() []byte {
	return bs.data
}

// getIndexes is a helper function that calculates the byte and bit index within the byte
// for the given bit position in the bitsream.
func getIndexes(index int) (byteInex, bitIndex int) {
	return index / bitSize, index % bitSize
}

// FromFile reads a file containing a list of numbers and returns a Bitstream.
func FromFile(filename string) (*BitStream, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if num, err := strconv.Atoi(line); err == nil {
			if num <= 0xff {
				data = append(data, byte(num))
			} else {
				data = append(data, byte(num>>8), byte(num&0xff))
			}
		} else {
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return NewBitStream(data), nil
}
