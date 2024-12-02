package main

import (
	"errors"
	"fmt"
	"lzw/decode"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	fileName, err := parseArgs()
	if err != nil {
		fmt.Println("arguments Error:", err)
		return
	}

	decodedBytes, err := decode.LZWDecode(fileName)
	if err != nil {
		fmt.Println("error decoding file:", err)
		return
	}
	fmt.Printf("decoded file: %s\n", fileName)

	err = writeToFile(fileName, decodedBytes)
	if err != nil {
		fmt.Println("failed to write to file:", err)
		return
	}
	fmt.Printf("decoded file written to: %s\n", fileName)

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

func writeToFile(fileName string, decodedBytes []byte) error {

	// get the name, and remove the .z extension
	base := filepath.Base(fileName)
	nameWithoutExt := strings.TrimSuffix(base, ".z")
	decodedFile := "decoded-" + nameWithoutExt

	file, err := os.Create(decodedFile)

	if err != nil {
		return errors.New("could not make a new file")
	}
	defer file.Close()

	_, err = file.Write(decodedBytes)
	if err != nil {
		return errors.New("could not write bytes to new file")
	}
	return nil
}
