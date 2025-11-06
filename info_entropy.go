/*
Provide structs to model Azul event metadata.
*/
package main

// Entropy structure.
type EventInfoEntropy struct {
	Overall    float64   `json:"overall"`
	BlockSize  int       `json:"block_size"`
	BlockCount int       `json:"block_count"`
	Blocks     []float64 `json:"blocks"`
}
