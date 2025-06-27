package oodle

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"log"
)

// ensureOodleLibrary checks if the Oodle library exists and downloads it if not.
// It calls log.Fatalf if the library cannot be made available, exiting all tests.
func ensureOodleLibraryLoaded() {
	if !IsLibExists() {
		fmt.Println("Oodle library not found. Attempting to download...")
		err := Download()
		if err != nil {
			if !IsLibExists() {
				log.Fatalf("Failed to download Oodle library: %v. Please ensure %s is available or downloadable.", err, libName)
			}
		}
		if !IsLibExists() {
			log.Fatalf("Oodle library was reportedly downloaded, but IsLibExists still reports it as missing. Check libName and download path.")
		}
		fmt.Printf("Oodle library downloaded successfully to: %s\n", getTempDllPath())
	} else {
		libPath, _ := resolveLibPath()
		fmt.Printf("Oodle library found at: %s\n", libPath)
	}

	_, err := loadLib()
	if err != nil {
		libPathForMsg, _ := resolveLibPath()
		if libPathForMsg == "" {
			libPathForMsg = fmt.Sprintf("one of the expected locations: current directory or %s", getTempDllPath())
		}
		log.Fatalf("Failed to load Oodle library from %s: %v. Ensure the library is correct for your system and architecture.", libPathForMsg, err)
	}
	fmt.Println("Oodle library loaded successfully for all tests.")
}

// TestMain will set up the Oodle library once for all tests in this package.
func TestMain(m *testing.M) {
	// Perform global setup, like ensuring Oodle DLL is available and loaded.
	ensureOodleLibraryLoaded()
	// Run the tests
	os.Exit(m.Run())
}

// generateCompressibleTestData creates a byte slice with repetitive patterns.
func generateCompressibleTestData(size int) []byte {
	if size == 0 {
		return []byte{}
	}
	pattern := []byte("RepeatThisPatternForGoodCompression! ")
	data := make([]byte, 0, size)
	for len(data) < size {
		data = append(data, pattern...)
	}
	return data[:size]
}

// TestAllCompressorsRoundTrip tests various compressors and levels with generated compressible data.
func TestAllCompressorsRoundTrip(t *testing.T) {
	// Generate sample data once
	sampleDataSize := 256 * 1024 // 256KB of compressible data
	sampleData := generateCompressibleTestData(sampleDataSize)
	if err := os.WriteFile("compressible_test_data.bin", sampleData, 0644); err != nil {
		t.Logf("Warning: failed to write test data to disk for inspection: %v", err)
	}


	t.Logf("Generated sample compressible data size: %d bytes", len(sampleData))

	compressorsToTest := []struct {
		Name       string
		Compressor int
	}{
		{"Kraken", CompressorKraken},
		{"Leviathan", CompressorLeviathan},
		{"Mermaid", CompressorMermaid},
		{"Selkie", CompressorSelkie},
		{"Hydra", CompressorHydra},
	}

	compressionLevelsToTest := []struct {
		Name  string
		Level int
	}{
		{"SuperFast", CompressionLevelSuperFast},
		{"Normal", CompressionLevelNormal},
		{"Optimal1", CompressionLevelOptimal1},
		{"Max", CompressionLevelMax},
		// Add HyperFast if desired, e.g. {"HyperFast", CompressionLevelHyperFast},
	}

	for _, compressorInfo := range compressorsToTest {
		t.Run(compressorInfo.Name, func(t *testing.T) {
			for _, levelInfo := range compressionLevelsToTest {
				// Capture range variables for subtests
				currentCompressorInfo := compressorInfo
				currentLevelInfo := levelInfo
				testName := fmt.Sprintf("Level_%s", currentLevelInfo.Name)

				t.Run(testName, func(t *testing.T) {
					t.Parallel() // Mark subtest as parallelizable if tests are independent
					t.Logf("Testing compressor: %s (%d) with level %s (%d)", currentCompressorInfo.Name, currentCompressorInfo.Compressor, currentLevelInfo.Name, currentLevelInfo.Level)

					compressedData, err := Compress(sampleData, currentCompressorInfo.Compressor, currentLevelInfo.Level)
					if err != nil {
						t.Errorf("Compress failed for %s with level %s: %v", currentCompressorInfo.Name, currentLevelInfo.Name, err)
						return
					}
					if len(compressedData) == 0 && len(sampleData) > 0 { // Compressing non-empty to empty is an error
						t.Errorf("Compress for %s with level %s returned empty data for non-empty input", currentCompressorInfo.Name, currentLevelInfo.Name)
						return
					}
					t.Logf("Compressed size with %s, level %s: %d bytes", currentCompressorInfo.Name, currentLevelInfo.Name, len(compressedData))

					decompressedData, err := Decompress(compressedData, int64(len(sampleData)))
					if err != nil {
						t.Errorf("Decompress failed for %s with level %s: %v", currentCompressorInfo.Name, currentLevelInfo.Name, err)
						return
					}

					if !bytes.Equal(sampleData, decompressedData) {
						t.Errorf("Decompressed data does not match original for %s with level %s.", currentCompressorInfo.Name, currentLevelInfo.Name)
					} else {
						t.Logf("Successfully compressed and decompressed with %s, level %s.", currentCompressorInfo.Name, currentLevelInfo.Name)
					}
				})
			}
		})
	}
}
