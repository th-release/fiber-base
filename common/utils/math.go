package utils

import (
	"fmt"
	"math"
	"strconv"
)

func ToFloat64(n int) float64 {
	return float64(n)
}

func ToInt(v interface{}) (int, error) {
	i, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("cannot convert %v to int", v)
	}
	return i, nil
}

func FloatToInt(f float64) int {
	return int(f)
}

func FloatToIntRound(f float64) int {
	return int(math.Round(f))
}

func MaxFloat64(values ...float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func MaxInt(values ...int) int {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// MinFloat64 returns the smaller of a and b, treating NaN as absent.
func MinFloat64(a, b float64) float64 {
	if math.IsNaN(a) {
		return b
	}
	if math.IsNaN(b) {
		return a
	}
	return math.Min(a, b)
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ToFixed rounds value to the given number of decimal places.
func ToFixed(value float64, decimals int) float64 {
	if decimals < 0 {
		decimals = 0
	}
	s := strconv.FormatFloat(value, 'f', decimals, 64)
	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return result
}

// PercentageDifference returns the absolute percentage difference of b relative to a.
func PercentageDifference(a, b float64) float64 {
	return math.Abs(a-b) / a * 100
}
