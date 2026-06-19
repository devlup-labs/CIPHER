package main

import (
	"strconv"
)

func GenerateLeafHash(
	fileID string,
	index int,
	length int,
	chunk []byte,
) string {

	combined := make([]byte, 0)

	// FileID
	combined = append(combined, []byte(fileID)...)

	// Chunk index
	combined = append(
		combined,
		[]byte(strconv.Itoa(index))...,
	)

	// Chunk length
	combined = append(
		combined,
		[]byte(strconv.Itoa(length))...,
	)

	// Raw chunk bytes
	combined = append(combined, chunk...)

	// Final leaf hash
	return GenerateHash(combined)
}