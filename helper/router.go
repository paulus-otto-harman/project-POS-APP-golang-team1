package helper

import (
	"strconv"
)

func Uint(param string) (uint, error) {
	i, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}
