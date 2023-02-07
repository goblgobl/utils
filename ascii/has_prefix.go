package ascii

func HasPrefixIgnoreCase(input string, prefix string) bool {
	if len(input) < len(prefix) {
		return false
	}

	for i := 0; i < len(prefix); i++ {
		a := prefix[i]
		b := input[i]
		if a == b {
			continue
		}
		if 'A' <= b && b <= 'Z' && a == b+32 {
			continue
		}

		if 'A' <= a && a <= 'Z' && b == a+32 {
			continue
		}

		return false
	}

	return true
}
