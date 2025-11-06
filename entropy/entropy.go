/*
Calculate Shannon's entropy over given byte buffers.
Supports calculating overall result or reporting as list of block size/count chunks.
*/
package entropy

import "math"

const MinBlockSize = 256

type Entropy struct {
	buf []byte
}

// Wrap a buffer for Shannon's entropy calculations
func New(buf []byte) (entropy *Entropy) {
	return &Entropy{buf}
}

// Value will calculate the entropy over all bytes
func (e *Entropy) Value() (answer float64) {
	var counts [256]int
	var px float64

	for _, b := range e.buf {
		counts[b]++
	}
	for i := 0; i < 256; i++ {
		if counts[i] == 0 {
			continue
		}
		px = float64(counts[i]) / float64(len(e.buf))
		answer -= px * math.Log2(px)
	}
	return answer
}

// BySize will calculate entropy over the bytes, split into
// chunks of the specified size (minimum of 256 byte chunks)
// the size and count of chunks are also returned
func (e *Entropy) BySize(size int) ([]float64, int, int) {
	if size < MinBlockSize {
		size = MinBlockSize
	}
	count := len(e.buf) / size
	ent := make([]float64, count)
	for i := 0; i < len(ent); i++ {
		start := i * size
		end := start + size
		ent[i] = New(e.buf[start:end]).Value()
	}
	return ent, size, count
}

// ByCount will calculate entropy over the bytes, split into
// equal sized chunks of the specified count (minimum of 256 byte chunks)
// the size and count of chunks are also returned
func (e *Entropy) ByCount(max_count int) ([]float64, int, int) {
	if max_count <= 0 {
		return e.BySize(0)
	}
	return e.BySize(len(e.buf) / max_count)
}
