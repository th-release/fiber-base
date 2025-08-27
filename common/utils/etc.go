package utils

import (
	"fmt"
	"math"
	"strconv"
)

func ToFloat64(num int) float64 {
	return float64(num)
}

func ToInt(v interface{}) (int, error) {
	i, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("value is not an int: %v", v)
	}
	return i, nil
}

func CalculateQuantity(
	price float64,
	usdt float64,
	leverage int,
	percent float64,
) float64 {
	investmentRatio := percent / 100
	investmentAmount := usdt * investmentRatio
	leveragedInvestment := investmentAmount * float64(leverage)
	return leveragedInvestment / price
}

func MaxFloat64(numbers ...float64) float64 {
	if len(numbers) == 0 {
		return 0
	}

	maxValue := numbers[0]
	for _, value := range numbers {
		if value > maxValue {
			maxValue = value
		}
	}

	return maxValue
}

func MaxInt(numbers ...int) int {
	if len(numbers) == 0 {
		return 0
	}

	maxValue := numbers[0]
	for _, value := range numbers {
		if value > maxValue {
			maxValue = value
		}
	}

	return maxValue
}

func MinFloat64(a, b float64) float64 {
	if math.IsNaN(a) && math.IsNaN(b) {
		return math.NaN()
	}
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

func FloatToInt(f float64) int {
	return int(f)
}

func FloatToIntWithRound(f float64) int {
	return int(math.Round(f))
}

func ToFixed(value float64, decimalPlaces int) float64 {
	if decimalPlaces < 0 {
		decimalPlaces = 0
	}

	format := fmt.Sprintf("%%.%df", decimalPlaces)

	formatted := fmt.Sprintf(format, value)

	result, err := strconv.ParseFloat(formatted, 64)
	if err != nil {
		return 0
	}

	return result
}

func PercentageDifference(a float64, b float64) float64 {
	diff := a - b
	abs := math.Abs(diff)
	return abs / a * 100
}
