package main

import (
	"fmt"
)

func main() {
	// Open the compressed file
	// Let x be 3 byte []byte with value 0xE2 0x80 0x9C

	x := []byte{0xE2, 0x80, 0x9C}

	// let us print the string value of this
	fmt.Println(string(x))
}
