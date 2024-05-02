package bitstream

// BitStreamWriter provides a way to write bits to a Bitstream.
type BitStreamWriter struct {
	bs   *BitStream // The Bitstream to write to
	data []byte     // The buffer for storing bits before flushing
	len  int        // The number of bits written to the buffer
}

// NewBitStreamWriter creates a new BitStreamWriter from the given Bitstream.
func NewBitStreamWriter(bs *BitStream) *BitStreamWriter {
	return &BitStreamWriter{bs: bs}
}

// Write writes the provided byte slice to the BitStreamWriter.
// Each byte is written as 8 individual bits.
func (w *BitStreamWriter) Write(bits []byte) error {
	for _, b := range bits {
		for i := msbIndex; i >= 0; i-- {
			bit := (b >> i) & 1
			w.data = append(w.data, bit)
			w.len++
		}
	}
	return nil
}

// WriteBit writes a single bit to the BitStreamWriter.
// It returns an error if the given bit value is not 0 or 1.
func (w *BitStreamWriter) WriteBit(bit byte) error {
	if bit != 0 && bit != 1 {
		return ErrInvalidBitValue
	}
	w.data = append(w.data, bit)
	w.len++
	return nil
}

// Flush writes the buffred bits to the Bitstream.
// It clears the buffer after writing.
func (w *BitStreamWriter) Flush() error {
	for _, bit := range w.data {
		err := w.bs.Append(bit)
		if err != nil {
			return err
		}
	}
	w.data = nil
	w.len = 0
	return nil
}
