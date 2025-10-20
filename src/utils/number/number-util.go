package numberUtil

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func RangeInt(min int, max int) int {
	return rand.Intn(max) + min
}

func Uint64ToString(number uint64) string {
	return strconv.FormatUint(number, 10)
}

func Int64ToString(number int64) string {
	return strconv.FormatInt(number, 10)
}

var StringToInt64 = func(str string) (int64, error) {
	number, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return number, nil
}

var StringToUint64 = func(str string) (uint64, error) {
	number, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func IntToString(number int) string {
	return strconv.Itoa(number)
}

func GetDecimalLength(number decimal.Decimal) int {
	str := number.String()
	parts := strings.Split(str, ".")
	if len(parts) == 2 {
		return len(parts[1])
	}
	return 0
}

func IsStringFloatValue(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}

func IsStringIntValue(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func ParseToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case uint64:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return num, nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", value)
	}
}
