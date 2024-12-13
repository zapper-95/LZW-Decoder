package encode

import (
	"os"
)

func initialiseMap(inputToEncoded map[string]uint16) {
	for i := 0; i < 256; i++ {
		inputToEncoded[string(i)] = uint16(i)
	}
}

func LZWEncode(fileName string) ([]byte, error) {
	// we want to be able to map an unbounded amount of bytes to codes between 0 and 2^12
	stringToCode := make(map[string]uint16)
	initialiseMap(inputToEncoded)

	codes, err := convertToCodes(fileName, stringToCode)

	// read in 12 bits

}

func convertToCodes(fileName string, stringToCode map[string]uint16) ([]uint16, error) {

	fp := 0
	rp := 1
	codes := make([]uint16, 0)

	for {
		file, err := os.Open(fileName)
		if err != nil {
			return nil, nil
		}

		buffer := make([]byte, 1)

		n, err := file.Read(buffer)
		if err != nil {
			return nil, nil
		}

		currentByte := append([]byte(nil), buffer...)

		for {
			file.Read(buffer)
			currentByte := append(currentByte, buffer...)
			_, ok := stringToCode[string(currentByte)]
			if !ok {
				// add to dictionary as new value
				stringToCode[string(currentByte)] = uint16(len(stringToCode))
			}
		}

	}
}
