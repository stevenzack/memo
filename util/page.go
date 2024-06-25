package util

func PageNum(total int64, size int) int {
	odd := total % int64(size)
	if odd == 0 {
		return int(total) / size
	}
	return int(total)/size + 1
}
