package main

import (
	"errors"
	"fmt"
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

	decodedString, err := decodeLZW(codes)

	if err != nil {
		fmt.Println("Error decoding file:", err)
	}

	err = writeToFile(fileName, decodedString)
	if err != nil {
		fmt.Println("Failed to write to file:", err)
	}
}

func writeToFile(fileName string, decodedString string) error {

	base := filepath.Base(fileName)
	nameWithoutExt := strings.TrimSuffix(base, ".z")
	decodedFile := "decoded-" + nameWithoutExt

	file, err := os.Create(decodedFile)

	if err != nil {
		return errors.New("Could not make new file")
	}
	defer file.Close()

	_, err = file.WriteString(decodedString)
	if err != nil {
		return errors.New("Could not write to new file")
	}
	return nil
}

func initialiseMap(codeToSymbol map[uint32]string) {
	for i := 0; i < 256; i++ {
		codeToSymbol[uint32(i)] = string(i)
	}

}

func decodeLZW(codes []uint32) (string, error) {

	codeToSymbol := make(map[uint32]string)
	initialiseMap(codeToSymbol)
	decodedSymbols := make([]string, len(codes))

	var prevSymbol string

	if len(codes) == 0 {
		return "", nil
	}

	for i, code := range codes {

		var currSymbol string

		// if the current code is in the map
		if symbol, ok := codeToSymbol[code]; ok {
			currSymbol = symbol
			// if the current code is the next to be added to the map
		} else if code == uint32(len(codeToSymbol)) && prevSymbol != "" {
			// utf-8 stores codes > 127 with multiple bits. Convert to rune to make sure I get the right char
			prevSymbolRunes := []rune(prevSymbol)
			currSymbol = prevSymbol + string(prevSymbolRunes[0])
		} else {
			return "", errors.New("Invalid code")
		}
		decodedSymbols[i] = currSymbol

		if prevSymbol != "" {

			// If there are no codes or the limit, we need to reset and add the first new character
			currSymbolRunes := []rune(currSymbol)
			newEntry := prevSymbol + string(currSymbolRunes[0])
			codeToSymbol[uint32(len(codeToSymbol))] = newEntry

			if len(codeToSymbol) == 4096 {
				fmt.Println(len(codeToSymbol))
				codeToSymbol = make(map[uint32]string)
				initialiseMap(codeToSymbol)
			}

		}

		prevSymbol = currSymbol

	}
	return strings.Join(decodedSymbols, ""), nil
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
