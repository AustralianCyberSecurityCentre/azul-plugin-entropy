package entropy

import (
	"reflect"
	"strings"
	"testing"
)

func TestEntropyBuffered(t *testing.T) {
	nulls := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	shrug := []byte{10, 88, 255, 13, 128, 77, 99, 123, 54}
	tables := []struct {
		input  []byte
		output float64
	}{
		{[]byte(""), 0.0},
		{[]byte("1223334444"), 1.8464393446710154},
		{[]byte("The quick brown fox jumps over the lazy dog"), 4.431965045349459},
		{nulls, 0.0},
		{shrug, 3.169925001442312},
		{[]byte(LargeBuffer), 4.380428799939244},
		// Length of below is specifically 307,618bytes and catches a potential out of bounds error
		{[]byte(strings.Repeat("AB", 153809)), 1},
	}

	for _, table := range tables {
		entropy := NewBuffered(uint64(len(table.input)), 800)
		entropy.AppendAndCalculateBufferedValues(table.input)
		tv, err := entropy.TotalValue()
		if err != nil {
			t.Errorf("error %v", err)
		}
		if tv != table.output {
			t.Errorf("Unexpected Entropy for: %v, expected %v got: %v", len(table.input), table.output, tv)
		}
	}
}

func TestEntropyBufferedChunks(t *testing.T) {
	tables := []struct {
		input       []byte
		count       int
		output      []float64
		outputSize  int
		outputCount int
	}{
		{[]byte(""), 0, []float64{}, 256, 0},           // too small to chunk
		{[]byte(""), 100, []float64{}, 256, 0},         // too small to chunk
		{[]byte("1223334444"), 1, []float64{}, 256, 0}, // too small to chunk
		{[]byte(LargeBuffer), 1, []float64{4.380428799939244}, 3876, 1},
		{[]byte(LargeBuffer), 5, []float64{4.443621850692178, 4.351689888387683, 4.292239846194779, 4.301192135704729, 4.241109953978476}, 775, 5},
		{[]byte(LargeBuffer), 10, []float64{4.399466412895255, 4.407784952821555, 4.237577608184258, 4.375320517593471, 4.261001802156862, 4.2480014235261665, 4.304663798108217, 4.229069528377996, 4.25618532350192, 4.1816184291151774}, 387, 10},
		{[]byte(LargeBuffer), 100, []float64{4.388541092008773, 4.3324519806257005, 4.348575276835077, 4.242764573097982, 4.176067779953571, 4.34031853435168, 4.26282964169927, 4.235566890498516, 4.231949341302883, 4.124254926105486, 4.370857649462027, 4.163746907414455, 4.223515182200197, 4.176315653309257, 4.112783102062334}, 256, 15},
	}
	for _, table := range tables {
		inLen := uint64(len(table.input))
		ent := NewBuffered(inLen, table.count)
		ent.AppendAndCalculateBufferedValues(table.input)
		entropy, size, count := ent.GetChunkEntropySizeAndCount()
		if size != table.outputSize {
			t.Errorf("Unexpected Output Size for: %v, got: %v", table.outputSize, size)
		}
		if count != table.outputCount {
			t.Errorf("Unexpected Output Count for: %v, got: %v", table.outputCount, count)
		}
		if !reflect.DeepEqual(entropy, table.output) {
			t.Errorf("Unexpected Entropy for: %v, got: %v", table.output, entropy)
		}
	}
}

// Test Entropy calculations still work when we append data in arbitrary byte slices.
func TestEntropyBufferedMultipleAppends(t *testing.T) {
	input := []byte(LargeBuffer)
	input_len := len(input)
	max_count := 100

	output_chunk_entropies := []float64{4.388541092008773, 4.3324519806257005, 4.348575276835077, 4.242764573097982, 4.176067779953571, 4.34031853435168, 4.26282964169927, 4.235566890498516, 4.231949341302883, 4.124254926105486, 4.370857649462027, 4.163746907414455, 4.223515182200197, 4.176315653309257, 4.112783102062334}
	outputSize := 256
	outputCount := 15
	outputTotalEntropy := 4.380428799939244

	for sliceSize := 1; sliceSize <= 3900; sliceSize++ {
		ent := NewBuffered(uint64(input_len), max_count)

		for i := 0; i <= input_len-1; i += sliceSize {
			end := i + sliceSize
			if i+sliceSize > input_len {
				end = input_len
			}
			ent.AppendAndCalculateBufferedValues(input[i:end])
		}

		entropy, size, count := ent.GetChunkEntropySizeAndCount()

		tv, err := ent.TotalValue()
		if err != nil {
			t.Errorf("error %v", err)
		}
		if tv != outputTotalEntropy {
			t.Errorf("SliceSize %d - Unexpected Entropy, expected %f got: %f", sliceSize, outputTotalEntropy, tv)
		}
		if size != outputSize {
			t.Errorf("SliceSize %d - Unexpected Output Size expected: %v, got: %v", sliceSize, outputSize, size)
		}
		if count != outputCount {
			t.Errorf("SliceSize %d - Unexpected Output Count expected: %v, got: %v", sliceSize, outputCount, count)
		}
		if !reflect.DeepEqual(entropy, output_chunk_entropies) {
			t.Errorf("SliceSize %d - Unexpected Entropy expected: %v, got: %v", sliceSize, output_chunk_entropies, entropy)
		}

	}
}

func BenchmarkEntropyBuffered(b *testing.B) {
	e := NewBuffered(uint64(len([]byte(LargeBuffer))), 100)
	e.AppendAndCalculateBufferedValues([]byte(LargeBuffer))
	for n := 0; n < b.N; n++ {
		_, err := e.TotalValue()
		if err != nil {
			panic(err)
		}
	}
}
