package util

import "testing"

func TestPageNum(t *testing.T) {
	n := PageNum(10, 5)
	if n != 2 {
		t.Error("n is not 2 , but ", n)
		return
	}
	n = PageNum(11, 5)
	if n != 3 {
		t.Error("n is not 3 , but ", n)
		return
	}
	n = PageNum(14, 5)
	if n != 3 {
		t.Error("n is not 3 , but ", n)
		return
	}
}
