package ascii

// https://www.openmymind.net/String-To-Integer-atoi-in-Go/

func Atoi(input string) (int, string) {
	var n int
	for i, b := range []byte(input) {
		b -= '0'
		if b > 9 {
			return n, input[i:]
		}
		n = n*10 + int(b)
	}
	return n, ""
}
