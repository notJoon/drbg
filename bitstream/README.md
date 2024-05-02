# Bitstream

The `bitstream` package provides a manipulating sequences of bits, offering a flexible way to handle bit-level data within bytes. This package is ideal for appliacations requiring direct bit manipulation such as compression algorithms, cryptography, and network protocols.

## Features

- **Create Bitstreams** from byte slices.
- **Read and Write** operations for individual bits.
- **Append bits**
- ** Stream bits** to a writer and reader.

## Usage

### Create a Bitstream

```go
func main() {
    data := []byte{0x0F, 0xF0}
    bs := bitstream.NewBitstream(data)

    // Print the length of bitstream
    fmt.Println("Length of Bitstream:", bs.Len())
}
```

### Writing to a Bitstream

```go
bs := bitstream.NewBitstream([]byte{})
writer := bitstream.NewBitStreamWriter(bs)

// Write a byte
writer.Write([]byte{0xAA})
writer.Flush()

// Write a single bit
writer.WriteBit(1)
writer.Flush()
```

### Reading from a Bitstream

```go
reader := bitstream.NewBitStreamReader(bs)

// Read a single byte
byteVal, err := reader.ReadByte()
if err != nil {
    log.Fatal(err)
}

// Read a single bit
bitVal, err := reader.ReadBit()
if err != nil {
    log.Fatal(err)
}
```
