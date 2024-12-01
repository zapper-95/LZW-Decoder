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

	decodedString, err := decodeLZW(codes)
	print(decodedString)
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
		codeToSymbol[uint32(i)] = string(rune(i))
	}

}

func decodeLZW(codes []uint32) (string, error) {

	codeToSymbol := make(map[uint32]string)
	initialiseMap(codeToSymbol)
	decodedSymbols := make([]string, 0)

	if len(codes) == 0 {
		return "", nil
	}

	firstCode := codes[0]

	prevSymbol, ok := codeToSymbol[firstCode]

	if !ok {
		return "", errors.New("First symbol not in map")
	}
	decodedSymbols = append(decodedSymbols, prevSymbol)

	for _, code := range codes[1:] {

		var currSymbol string

		// if the current code is in the map
		if symbol, ok := codeToSymbol[code]; ok {
			currSymbol = symbol
			// if the current code is the next to be added to the map
		} else if code == uint32(len(codeToSymbol)) {
			// utf-8 stores codes > 127 with multiple bits. Convert to rune to make sure I get the right char

			currSymbol = prevSymbol + string(([]rune(prevSymbol))[0])
		} else {
			return "", errors.New("Invalid code")
		}
		decodedSymbols = append(decodedSymbols, currSymbol)

		// If there are no codes or the limit, we need to reset and add the first new character
		if prevSymbol == "â" {
			fmt.Println(prevSymbol)
			fmt.Println(prevSymbol)
			break
		}
		newEntry := prevSymbol + string([]rune(currSymbol)[0])
		//newEntry := append(prevSymbol, currSymbol[0])
		codeToSymbol[uint32(len(codeToSymbol))] = newEntry

		if len(codeToSymbol) == 4096 {
			codeToSymbol = make(map[uint32]string)
			initialiseMap(codeToSymbol)
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
