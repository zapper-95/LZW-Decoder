package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
)

func main() {
	fileName, err := parseArgs()

	if err != nil {
		fmt.Println("Arguments Error:", err)
		return
	}

	codes, err := getCodes(fileName)

	if err != nil {
		fmt.Println("File splitting error:", err)
	}

	decodedString := decodeLZW(codes)

	fmt.Println(decodedString)

}

func initialiseMap(codeToSymbol map[uint32]string) {
	for i := 0; i < 256; i++ {
		codeToSymbol[uint32(i)] = string(i)
	}

}

func decodeLZW(codes []uint32) string {
	codeToSymbol := make(map[uint32]string)
	initialiseMap(codeToSymbol)

	decodedSymbols := make([]string, len(codes))

	// first code must be a number between 0 and 255
	decodedSymbols[0] = string(codes[0])

	for i, code := range codes {
		if i == 0 {
			continue
		}

		// new symbol to be added to the map
		var newMapSymbol string

		prevSymbol := decodedSymbols[i-1]

		// if the current symbol is in the map
		if currSymbol, ok := codeToSymbol[code]; ok {
			decodedSymbols[i] = currSymbol

			newMapSymbol = prevSymbol + currSymbol[0:1]

		} else {
			newMapSymbol = prevSymbol + prevSymbol[0:1]
			decodedSymbols[i] = newMapSymbol
		}

		// Add the new map symbol to the map
		if newMapSymbol != "" {

			// all possible codes have been used, so reset map
			if len(codeToSymbol) >= int(math.Pow(2, 12))-256 {
				codeToSymbol := make(map[uint32]string)
				initialiseMap(codeToSymbol)
			}
			symbolsCount := uint32(len(codeToSymbol))
			codeToSymbol[symbolsCount] = newMapSymbol
		}

	}
	return strings.Join(decodedSymbols, "")
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

func getCodes(fileName string) ([]uint32, error) {
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

	// Bitshift the last code right by 4 if the number of codes is odd
	if len(codes)%2 != 0 {
		codes[len(codes)-1] = codes[len(codes)-1] >> 4
	}

	// Number of bits not divisble by 12
	if collectedBitsSize%2 != 0 {
		return nil, errors.New("not a multiple of 12 bits")
	}

	return codes, nil
}
