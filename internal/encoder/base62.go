package encoder

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(n uint64) string {
	if n == 0 {
		return "0"
	}

	b := make([]byte, 11)
	i := len(b)

	for n > 0 {
		i--
		b[i] = alphabet[n%62]
		n /= 62
	}

	return string(b[i:])
}