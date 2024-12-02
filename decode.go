package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

	decodedBytes, err := decodeLZW(codes)
	print(decodedBytes)
	if err != nil {
		fmt.Println("Error decoding file:", err)
	}

	err = writeToFile(fileName, decodedBytes)
	if err != nil {
		fmt.Println("Failed to write to file:", err)
	}
}

func writeToFile(fileName string, decodedBytes []byte) error {

	fmt.Printf("%d", decodedBytes)
	base := filepath.Base(fileName)
	nameWithoutExt := strings.TrimSuffix(base, ".z")
	decodedFile := "decoded-" + nameWithoutExt

	file, err := os.Create(decodedFile)

	if err != nil {
		return errors.New("Could not make new file")
	}
	defer file.Close()

	_, err = file.Write(decodedBytes)
	if err != nil {
		return errors.New("Could not write to new file")
	}
	return nil
}

// initialiseMap initializes the codeToSymbol map with codes from 0 to 255.
func initialiseMap(codeToSymbol map[uint32][]byte) {
	for i := 0; i < 256; i++ {
		codeToSymbol[uint32(i)] = []byte{byte(i)}
	}
}

// decodeLZW decodes a slice of LZW codes into the original byte sequence.
func decodeLZW(codes []uint32) ([]byte, error) {
	// Initialize the codeToSymbol map with a range of 256.
	codeToSymbol := make(map[uint32][]byte)
	initialiseMap(codeToSymbol)
	decodedSymbols := make([]byte, 0)

	if len(codes) == 0 {
		return nil, nil
	}

	firstCode := codes[0]

	prevSymbol, ok := codeToSymbol[firstCode]
	if !ok {
		return nil, errors.New("first symbol not in map")
	}

	decodedSymbols = append(decodedSymbols, prevSymbol...)

	for _, code := range codes[1:] {
		var currSymbol []byte

		if symbol, ok := codeToSymbol[code]; ok {
			// Current code is in the map.
			currSymbol = symbol
		} else if code == uint32(len(codeToSymbol)) {
			// Current code is the next to be added to the map.
			// This handles the special case in LZW where the code is not yet in the map.
			currSymbol = append(append([]byte(nil), prevSymbol...), prevSymbol[0])
		} else {
			return nil, errors.New("invalid code encountered during decoding")
		}

		// Append the entire current symbol to decodedSymbols.
		decodedSymbols = append(decodedSymbols, currSymbol...)

		// Add new entry to the map: previous symbol + first byte of current symbol.
		newEntry := append(append([]byte(nil), prevSymbol...), currSymbol[0])
		codeToSymbol[uint32(len(codeToSymbol))] = newEntry

		// Reset the map if it reaches the maximum size (e.g., 4096 entries).
		if len(codeToSymbol) >= 4096 {
			codeToSymbol = make(map[uint32][]byte)
			initialiseMap(codeToSymbol)
		}

		// Update prevSymbol for the next iteration.
		prevSymbol = currSymbol
	}

	return decodedSymbols, nil
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

	buffer := make([]byte, 3)

	for {
		// Reads in the current byte into the buffer
		n, err := file.Read(buffer)

		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if n >= 2 {
			// bytes are int8, so need to cast to allow bitshifting
			firstCode := (uint32(buffer[0]) << 4) | (uint32(buffer[1]) >> 4)
			codes = append(codes, firstCode)
		} else {
			return nil, errors.New("Invalid number of bytes")
		}

		if n == 3 {
			secondCode := ((uint32(buffer[1]) & 0x0F) << 8) | uint32(buffer[2])
			codes = append(codes, secondCode)
		}

	}
	return codes, nil
}
