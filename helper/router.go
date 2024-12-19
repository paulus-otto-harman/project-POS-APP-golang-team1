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

func Float(param string) (float64, error) {
	f, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}
