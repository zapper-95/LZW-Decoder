package decode

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const initChars uint16 = 256

func LZWDecode(fileName string) ([]byte, error) {

	codes, err := getCodes(fileName)
	if err != nil {
		fmt.Println("error splitting files bits:", err)
		return nil, err
	}

	decodedBytes, err := lzwDecodeBytes(codes)

	if err != nil {
		fmt.Println("error decoding file:", err)
		return nil, err
	}

	return decodedBytes, nil

}

// mapping from ints to corresponding bytes
func initialiseMap(codeToSymbol [][]byte) {
	for i := 0; i < int(initChars); i++ {
		codeToSymbol[uint16(i)] = []byte{byte(i)}
	}
}

// uses codes to decode a sequence of bytes
func lzwDecodeBytes(codes []uint16) ([]byte, error) {

	codeToByte := make([][]byte, 4096)
	// intialises the slice with the bytes of the first 256 symbols
	initialiseMap(codeToByte)
	nextCode := initChars

	decodedBytes := make([]byte, 0)

	if len(codes) == 0 {
		return nil, nil
	}

	// handle the first code separately which is guarenteed to be in the slice
	firstCode := codes[0]

	prevSymbol := codeToByte[firstCode]

	decodedBytes = append(decodedBytes, prevSymbol...)

	for _, code := range codes[1:] {
		var currSymbol []byte

		if code < nextCode {
			currSymbol = codeToByte[code]
		} else if code == nextCode {
			// this handles the special case, where the code is not yet in the slice
			// this only occurs when the encoder uses the previous code as the next symbol

			currSymbol = append(append([]byte(nil), prevSymbol...), prevSymbol[0])
		} else {
			return nil, errors.New("invalid code encountered during decoding")
		}

		// append the current symbol to the decoded bytes
		decodedBytes = append(decodedBytes, currSymbol...)

		// add the new entry to the slice
		newEntry := append(append([]byte(nil), prevSymbol...), currSymbol[0])
		codeToByte[nextCode] = newEntry

		nextCode += 1
		// reset the next code in the slice if it reaches the maximum size of 2^12
		if nextCode >= 4096 {
			nextCode = initChars
		}

		prevSymbol = currSymbol
	}

	return decodedBytes, nil
}

func getCodes(fileName string) ([]uint16, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("could not open the file")
	}
	defer file.Close()

	var codes []uint16

	// read in the bytes in chunks of 3
	buffer := make([]byte, 3)

	for {
		n, err := file.Read(buffer)

		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if n >= 2 {
			// first code is the 8 bits of the first byte, and the first 4 bits of the second byte
			firstCode := (uint16(buffer[0]) << 4) | (uint16(buffer[1]) >> 4)
			codes = append(codes, firstCode)
		} else {
			return nil, errors.New("too few number of bytes")
		}

		if n == 3 {
			// second code is the last 4 bits of the second byte and the 8 bits of the third byte
			secondCode := ((uint16(buffer[1]) & 0x0F) << 8) | uint16(buffer[2])
			codes = append(codes, secondCode)
		}

	}
	return codes, nil
}
