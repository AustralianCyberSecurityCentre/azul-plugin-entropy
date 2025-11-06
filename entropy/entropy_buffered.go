/*
Calculate Shannon's entropy over given byte buffers.
Supports calculating overall result or reporting as list of block size/count chunks.
*/
package entropy

import (
	"fmt"
	"math"
)

const BufferedMinBlockSize = 256

// Struct that buffers a Binaries Entropy and continually increments the counts calculating the Shannon's entropy
// The chunkEntropies holds the Entropy of all the chunks of a fixed size.
// The number of bytes in one chunk is calculated based on the total length of the file and the max_block_count which
// is the maximum number of chunks provided in constructor.
type EntropyBuffered struct {
	contentLength       uint64
	actualContentLength uint64
	totalCount          [256]int
	// Chunk Entropy Constants
	size  int
	count int
	// Chunk Entropy Tracking
	chunkEntropies    []float64
	currentChunkCount [256]int
	chunkSizeSoFar    int
	chunkCountIdx     int
}

// Creates a new BufferedEntropy for a binary that is contentLength bytes long and can have at most max_block_count
// entropy chunk blocks.
func NewBuffered(contentLength uint64, max_block_count int) (entropyBuffered *EntropyBuffered) {
	size, count := calcSizeAndCount(max_block_count, contentLength)
	return &EntropyBuffered{
		contentLength:       contentLength,
		actualContentLength: 0,
		size:                size,
		count:               count,
		chunkEntropies:      make([]float64, count),
		chunkCountIdx:       0,
	}
}

// Appends new data to the BufferedEntropy adding to the chunked and total entropy counts.
// If enough data for one or more chunks to be calculated is provided it calculates the entropy for the chunk(s).
func (eb *EntropyBuffered) AppendAndCalculateBufferedValues(buf []byte) {
	for _, b := range buf {
		// Increment who file entropy counter
		eb.totalCount[b]++
		eb.actualContentLength += 1

		// Increment chunk counter and length of chunk
		eb.currentChunkCount[b]++
		eb.chunkSizeSoFar += 1

		// If chunk has hit the max chunk size calculate the entropy for the chunk and clear out chunk counters.
		if eb.chunkSizeSoFar == eb.size {
			eb.calcChunkValueAndClearCount(eb.chunkSizeSoFar)
			eb.chunkSizeSoFar = 0
		}
	}
}

// Calculate and return the total Entropy of all bytes provided to the EntropyBuffer.
// Expected to be called once whole file has been appended to the buffer.
func (eb *EntropyBuffered) TotalValue() (float64, error) {
	if eb.actualContentLength != eb.contentLength {
		return 0, fmt.Errorf("expected %d bytes, but got %d bytes", eb.contentLength, eb.actualContentLength)
	}
	return calculateEntropy(eb.totalCount, eb.contentLength), nil
}

// Calculate the entropy for the current chunk and clear the current chunks counts.
// If the max_count has been reached the remaining data is discarded.
// This occurs if the max_count and content length have a wide gap (refer to readme.md)
func (eb *EntropyBuffered) calcChunkValueAndClearCount(chunkLength int) {
	// Discard left over data that couldn't fit in any blocks.
	if eb.chunkCountIdx >= eb.count {
		return
	}
	entropy := calculateEntropy(eb.currentChunkCount, uint64(chunkLength))
	eb.chunkEntropies[eb.chunkCountIdx] = entropy
	eb.chunkCountIdx += 1

	for i := 0; i < 256; i++ {
		eb.currentChunkCount[i] = 0
	}
}

// Calculates the Shannon's Entropy for the provided count and provided bytes.
// Inputs are the entropy counts, bufferLength number of bytes used to generate the count.
func calculateEntropy(counts [256]int, bufferLength uint64) (answer float64) {
	var px float64

	for i := 0; i < 256; i++ {
		if counts[i] == 0 {
			continue
		}
		px = float64(counts[i]) / float64(bufferLength)
		answer -= px * math.Log2(px)
	}
	return answer
}

// Calculates the number of blocks to chunk a file of a provided length up to the maximum number of blocks.
// And with a minimum size of 256bytes per block.
func calcSizeAndCount(max_count int, contentLength uint64) (int, int) {
	size := 0
	if max_count != 0 {
		size = int(contentLength / uint64(max_count))
	}
	count := max_count
	if size < BufferedMinBlockSize {
		size = BufferedMinBlockSize
		count = int(contentLength / uint64(size))
	}
	return size, count
}

// Get all the entropies for the file chunks, and provides the bytes in each chunk as well as the total count of
// entropies.
func (eb *EntropyBuffered) GetChunkEntropySizeAndCount() ([]float64, int, int) {
	return eb.chunkEntropies, eb.size, eb.count
}
