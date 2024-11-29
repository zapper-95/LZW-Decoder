# LZW-Decryption
1. Split the input.txt.z files into groups of 12 bits, which are stored inside a slice
    * Read in the file using command line arguments
    * Iterate through each byte
        * Add the current code to the buffer by bitshifting left the current buffer and applying a logical OR
        * If the length of this buffer >= 12, take the first 12
            * Do this by bitshifting right by len(buffer) - 12
            * Apend this to a slice of codes
            * Reduce buffer to only the remaing bits. We can do this using a mask of the len(buffer)-12 least significant bits and applying OR operation
        * If the number of codes is odd, we need to shift the last code by 4 to the right, since it has been padded.
    *
2. Pass this slice to a function, which decodes this slice
    * A map is made to map decimal values to the 256 ascii characters
    * Iterate through each code
        * If the code is in the map, output the code
            * Add old_code + new_code[0] to the map at the next space
        * Else the code must be made using old_code + new_code[0]. So the new_code = old_code + old_code[0]. 
            Output this, and add the same to the map
        * If the size of the map ever becomes larger than 2^12, we remove all indexes greater than 255.
    Return the output  


    Questiions:
    Read the whole file in vs read line by line?

    Read whole thing in assumes we have enough memory to do - might not for large files

    100100000000
    
    1100101

    What if there is not a multiple of 12 bits? Need to check that any collected bits is not less than 8 of them.