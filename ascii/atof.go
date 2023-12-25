package ascii

import (
	"strconv"

	"src.goblgobl.com/utils"
)

// https://www.openmymind.net/String-To-Integer-atoi-in-Go/

func Atof(input string) (float64, string) {
	dot := false
	i := 0
	for _, b := range utils.S2B(input) {
		i += 1
		if b < '0' || b > '9' {
			if b == '.' && dot == false {
				dot = true
				continue
			}
			break
		}
	}

	end := i - 1
	if end == len(input)-1 {
		end = i
	}
	flt, err := strconv.ParseFloat(input[:end], 64)
	if err != nil {
		return 0, input
	}
	if end == len(input)-1 {
		return flt, ""
	}

	return flt, input[end:]
}
