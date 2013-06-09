package util

func Hash(s string) int64 {
	var hash int64
	for _, rune := range s {
		hash = (hash * 397) ^ int64(rune)
	}
	return hash
}