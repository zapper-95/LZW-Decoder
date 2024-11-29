package main

import (
	"fmt"
	"os"
)

func main() {
	// Open the compressed file

	if len(os.Args) < 2 {
		fmt.Println("Please run: go run main.go <filename>")
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file byte-by-byte
	buffer := make([]byte, 1)
	bitBuffer := uint32(0) // Holds the bits read so far
	bitsInBuffer := 0      // Number of bits currently in the buffer

	for {
		// Read one byte
		n, err := file.Read(buffer)
		if n == 0 || err != nil {
			break // End of file or error
		}

		// Add the byte to the bit buffer
		bitBuffer = (bitBuffer << 8) | uint32(buffer[0])
		bitsInBuffer += 8

		// Process bits from the buffer (example: extracting 12-bit codes)
		for bitsInBuffer >= 12 {
			// Extract the top 12 bits as a code
			code := (bitBuffer >> (bitsInBuffer - 12)) & 0xFFF
			bitsInBuffer -= 12

			// Print the 12-bit code
			fmt.Printf("Code: %d\n", code)
		}
	}

	// If there are leftover bits in the buffer, handle them (optional)
	if bitsInBuffer > 0 {
		fmt.Printf("Remaining bits: %b\n", bitBuffer&((1<<bitsInBuffer)-1))
	}

}
