package util

func SubAbs(a, b uint16) uint16 {
	if a < b {
		return 0
	}
	return a - b
}
