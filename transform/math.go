package transform

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
)

// ParsePriceAmountFloat 价格和数量string转float64
func ParsePriceAmountFloat(data []string) (price float64, amount float64, err error) {
	price, err = strconv.ParseFloat(data[0], 64)
	if err != nil {
		return
	}
	amount, err = strconv.ParseFloat(data[1], 64)
	if err != nil {
		return
	}
	return
}

// XToString 未知类型转string
func XToString(value interface{}) (res string) {
	switch v := value.(type) {
	case int:
		res = strconv.Itoa(v)
	case int8:
		res = strconv.FormatInt(int64(v), 10)
	case int16:
		res = strconv.FormatInt(int64(v), 10)
	case int32:
		res = strconv.FormatInt(int64(v), 10)
	case int64:
		res = strconv.FormatInt(v, 10)
	case float64:
		res = strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		res = strconv.FormatFloat(float64(v), 'f', -1, 64)
	case string:
		res = v
	}
	return
}

// StringToX string转指定类型
func StringToX[T int | int64 | float64 | string](v string) (value interface{}) {
	var res interface{}
	var a T
	res = a
	switch res.(type) {
	case int:
		value, _ = strconv.Atoi(v)
	case int64:
		value, _ = strconv.ParseInt(v, 10, 64)
	case float64:
		value, _ = strconv.ParseFloat(v, 64)
	case string:
		value = v
	default:
		_ = errors.New("not support type")
	}
	return
}

func ToFloat64(v interface{}) float64 {
	if v == nil {
		return 0.0
	}

	switch v.(type) {
	case float64:
		return v.(float64)
	case string:
		vStr := v.(string)
		vF, _ := strconv.ParseFloat(vStr, 64)
		return vF
	default:
		return 0
	}
}

func ToInt(v interface{}) int {
	if v == nil {
		return 0
	}

	switch v.(type) {
	case string:
		vStr := v.(string)
		vInt, _ := strconv.Atoi(vStr)
		return vInt
	case int:
		return v.(int)
	case float64:
		vF := v.(float64)
		return int(vF)
	default:
		return 0
	}
}

func ToUint64(v interface{}) uint64 {
	if v == nil {
		return 0
	}

	switch v.(type) {
	case int:
		return uint64(v.(int))
	case float64:
		return uint64((v.(float64)))
	case string:
		uV, _ := strconv.ParseUint(v.(string), 10, 64)
		return uV
	default:
		return 0
	}
}

func ToInt64(v interface{}) int64 {
	if v == nil {
		return 0
	}

	switch v.(type) {
	case float64:
		return int64(v.(float64))
	default:
		vv := fmt.Sprint(v)

		if vv == "" {
			return 0
		}

		vvv, err := strconv.ParseInt(vv, 0, 64)
		if err != nil {
			return 0
		}

		return vvv
	}
}

func FloatToString(v float64, precision int) string {
	return fmt.Sprint(FloatToFixed(v, precision))
}

func FloatToFixed(v float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Round(v*p) / p
}

func ValuesToJson(v url.Values) ([]byte, error) {
	parammap := make(map[string]interface{})
	for k, vv := range v {
		if len(vv) == 1 {
			parammap[k] = vv[0]
		} else {
			parammap[k] = vv
		}
	}
	return json.Marshal(parammap)
}

func Min[T int64 | int | float64 | float32 | string](a ...T) (b T) {
	b = a[0]
	for _, i := range a[1:] {
		if i < b {
			b = i
		}
	}
	return
}

func findFunc[T comparable](a []T, v T) int {
	for i, e := range a {
		if e == v {
			return i
		}
	}
	return -1
}

func Max[T int64 | int | float64 | float32 | string](a ...T) (b T) {
	b = a[0]
	for _, i := range a[1:] {
		if i > b {
			b = i
		}
	}
	return
}

func Str2Int64(v string) (int64, error) {
	value, err := strconv.ParseInt(v, 10, 64)
	return value, err
}

func Str2Int32(v string) (int32, error) {
	value, err := strconv.ParseInt(v, 10, 32)
	return int32(value), err
}

func Str2Float64(v string) (float64, error) {
	value, err := strconv.ParseFloat(v, 64)
	return value, err
}

func FormatFloatE(f float64, r float64) string {
	if math.Abs(f) < r {
		s := strconv.FormatFloat(f, 'e', -1, 64)
		return strings.Replace(s, "e-0", "e-", 1)
	} else {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
}
