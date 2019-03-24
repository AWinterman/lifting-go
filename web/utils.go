package main

import (
	"strconv"
)

func parseInt(value string) (int, error) {
	if len(value) < 1 {
		return 0, nil
	}
	i, errs := strconv.Atoi(value)

	return i, errs
}
