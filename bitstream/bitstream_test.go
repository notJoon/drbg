package bitstream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBitstream(t *testing.T) {
	data := []byte{0xAB, 0xCD, 0xEF}
	bs := NewBitstream(data)
	assert.NotNil(t, bs)
	assert.Equal(t, 24, bs.Len())
}

func TestBitstreamLen(t *testing.T) {
	data := []byte{0xAB, 0xCD, 0xEF}
	bs := NewBitstream(data)
	assert.Equal(t, 24, bs.Len())
}

func TestBitstreamBit(t *testing.T) {
	data := []byte{0xAB, 0xCD, 0xEF}
	bs := NewBitstream(data)
	bit, err := bs.Bit(0)
	assert.NoError(t, err)
	assert.Equal(t, byte(1), bit)
	bit, err = bs.Bit(23)
	assert.NoError(t, err)
	assert.Equal(t, byte(1), bit)
}

func TestBitstreamSetBit(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00}
	bs := NewBitstream(data)
	err := bs.SetBit(0, 1)
	assert.NoError(t, err)
	err = bs.SetBit(23, 1)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x80, 0x00, 0x01}, bs.Bytes())
}

func TestBitstreamAppend(t *testing.T) {
	bs := NewBitstream([]byte{})
	err := bs.Append(1)
	assert.NoError(t, err)
	err = bs.Append(0)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x80}, bs.Bytes())
}

func TestBitstreamBytes(t *testing.T) {
	data := []byte{0xAB, 0xCD, 0xEF}
	bs := NewBitstream(data)
	assert.Equal(t, data, bs.Bytes())
}
