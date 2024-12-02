package decode

import (
	"errors"
	"fmt"
	"io"
	"os"
)

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
func initialiseMap(codeToSymbol map[uint32][]byte) {
	for i := 0; i < 256; i++ {
		codeToSymbol[uint32(i)] = []byte{byte(i)}
	}
}

// uses codes to decode a sequence of bytes
func lzwDecodeBytes(codes []uint32) ([]byte, error) {

	codeToByte := make(map[uint32][]byte)

	// intialises the map with the bytes of the first 256 symbols
	initialiseMap(codeToByte)

	decodedBytes := make([]byte, 0)

	if len(codes) == 0 {
		return nil, nil
	}

	// handle the first code separately which is guarenteed to be in the map
	firstCode := codes[0]

	prevSymbol, ok := codeToByte[firstCode]
	if !ok {
		return nil, errors.New("first symbol not in map")
	}

	decodedBytes = append(decodedBytes, prevSymbol...)

	for _, code := range codes[1:] {
		var currSymbol []byte

		if symbol, ok := codeToByte[code]; ok {
			currSymbol = symbol
		} else if code == uint32(len(codeToByte)) {
			// this handles the special case, where the code is not yet in the map
			// this only occurs when the encoder uses the previous code as the next symbol

			currSymbol = append(append([]byte(nil), prevSymbol...), prevSymbol[0])
		} else {
			return nil, errors.New("invalid code encountered during decoding")
		}

		// append the current symbol to the decoded bytes
		decodedBytes = append(decodedBytes, currSymbol...)

		// add the new entry to the map
		newEntry := append(append([]byte(nil), prevSymbol...), currSymbol[0])
		codeToByte[uint32(len(codeToByte))] = newEntry

		// reset the map if it reaches the maximum size of 2^12
		if len(codeToByte) >= 4096 {
			codeToByte = make(map[uint32][]byte)
			initialiseMap(codeToByte)
		}

		prevSymbol = currSymbol
	}

	return decodedBytes, nil
}

func getCodes(fileName string) ([]uint32, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("could not open the file")
	}
	defer file.Close()

	var codes []uint32

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
			firstCode := (uint32(buffer[0]) << 4) | (uint32(buffer[1]) >> 4)
			codes = append(codes, firstCode)
		} else {
			return nil, errors.New("too few number of bytes")
		}

		if n == 3 {
			// second code is the last 4 bits of the second byte and the 8 bits of the third byte
			secondCode := ((uint32(buffer[1]) & 0x0F) << 8) | uint32(buffer[2])
			codes = append(codes, secondCode)
		}

	}
	return codes, nil
}
