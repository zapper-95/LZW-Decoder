# LZW-Decoder
## Summary
LZW Decoder for 12 bit fixed width codes written in Go. The dictionary is initialised with the first 256 entries of possible bytes.

Run from main with `go run main.go [path to input file]`. The output file will be written to a file inside the main directory.


## Approach
My approach is split into two main parts: extracting the binary data into blocks of 12 bits and then iteratively building a map and storing decoded bytes. 

Each 12 bits is read, decoded, and then added to the decoded bytes. A mapping from codes to bytes is also updated. Once this mapping is full (i.e of size 2^12), it is set to be only be the first 256 again.
