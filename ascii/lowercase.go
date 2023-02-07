package ascii

import "unsafe"

// https://www.openmymind.net/ASCII_String_To_Lowercase_in_Go/
func Lowercase(input string) string {
	l := len(input)
	for i := 0; i < l; i++ {
		c := input[i]
		if 'A' <= c && c <= 'Z' {
			// We've found an uppercase character, we'll need to convert this string
			lower := make([]byte, l)

			// copy everything we've skipped over up to this point
			copy(lower, input[:i])

			// our current character needs to be uppercase (it's the reason we're
			// in this branch)
			lower[i] = c + 32

			// now iterate over the rest of the input, from where we are, knowing that
			// we goneed to copy/lower case into our lowercase strinr
			for i += 1; i < l; i++ {
				c := input[i]
				if 'A' <= c && c <= 'Z' {
					c += 32
				}
				lower[i] = c
			}
			// if you think this is unfair, note that strings.Builder
			// does the exact same thing
			return *(*string)(unsafe.Pointer(&lower))
		}
	}

	// input was already lowercase, return it as-is
	return input
}
