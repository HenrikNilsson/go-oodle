package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/new-world-tools/go-oodle"
)

func main() {
	fmt.Println("Starting Oodle test application...")
	// The go-oodle library itself will handle ensuring the DLL is loaded
	// when its functions like Compress/Decompress are called, using its
	// own internal logic (including download with embedded fallbacks).

	originalData := []byte("This is some sample data for testing the Oodle library from an external application. It should compress and decompress successfully.")
	fmt.Printf("Original data: \"%s\"\n", string(originalData))
	fmt.Printf("Original size: %d bytes\n\n", len(originalData))

	compressor := oodle.CompressorKraken
	level := oodle.CompressionLevelNormal

	fmt.Printf("Attempting to compress with Compressor: %d, Level: %d\n", compressor, level)
	compressedData, err := oodle.Compress(originalData, compressor, level)
	if err != nil {
		log.Fatalf("Error during compression: %v", err)
	}
	fmt.Printf("Compressed data size: %d bytes\n", len(compressedData))

	if len(compressedData) == 0 && len(originalData) > 0 {
		log.Fatalf("Compression returned empty data for non-empty input!")
	}
	if float64(len(compressedData)) >= float64(len(originalData))*0.9 { // Check if compression is significant
		fmt.Println("Note: Compressed data is not significantly smaller than original. This can happen for short, random, or already compressed data.")
	}
	fmt.Println("Compression successful.")

	fmt.Println("\nAttempting to decompress...")
	decompressedData, err := oodle.Decompress(compressedData, int64(len(originalData)))
	if err != nil {
		log.Fatalf("Error during decompression: %v", err)
	}
	fmt.Printf("Decompressed data size: %d bytes\n", len(decompressedData))
	fmt.Println("Decompression successful.")

	fmt.Println("\nVerifying data...")
	if bytes.Equal(originalData, decompressedData) {
		fmt.Println("SUCCESS: Decompressed data matches original data!")
	} else {
		log.Fatalf("FAILURE: Decompressed data does NOT match original data.\nOriginal: %s\nDecompressed: %s", string(originalData), string(decompressedData))
	}
}
