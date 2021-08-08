package main

import (
	"strconv"
)

type metricValueFunc func(string) (float64, error)

func fromNumeric(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}
