package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	fileName, err := parseArgs()

	if err != nil {
		fmt.Println("Arguments Error:", err)
		return
	}

	codes, err := splitInput(fileName)
	fmt.Println(codes)
	if err != nil {
		fmt.Println("File splitting error:", err)
	}

}

func parseArgs() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("needs a filename argument")
	}
	if len(os.Args) > 2 {
		return "", errors.New("only one argument required for filename")
	}

	return os.Args[1], nil

}

func splitInput(fileName string) ([]uint32, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("could not open the file")
	}
	defer file.Close()

	// A code is up to 12 bits long
	var codes []uint32

	buffer := make([]byte, 1)

	// The current bits collected is up to 11+8=19
	collectedBits := uint32(0)
	collectedBitsSize := 0

	for {
		// Reads in the current byte into the buffer
		n, err := file.Read(buffer)

		if n == 0 || err != nil {
			break
		}

		// shift the bits left and apply a logical or with the next 8 bits
		collectedBits = (collectedBits << 8) | uint32(buffer[0])
		collectedBitsSize += 8

		if collectedBitsSize >= 12 {
			// add the first 12 bits as a new code
			codes = append(codes, collectedBits>>(collectedBitsSize-12))

			// set the collected bits to be only the remaining unused bits
			collectedBits = collectedBits & ((1 << (collectedBitsSize - 12)) - 1)
			collectedBitsSize -= 12
		}

	}

	// Number of bits not divisble by 12
	if collectedBitsSize != 0 {
		return nil, errors.New("not a multiple of 12 bits")
	}

	return codes, nil
}
