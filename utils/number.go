package utils

import (
	"fmt"
	"math"
	"strings"
)

func FormatNumber(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return FormatNumber(s[:n-3]) + "," + s[n-3:]
}

/**
 * cnNiceNumber 「万亿」单位
 * - 小于 1 万：直接展示
 * - 小于 1 亿：单位为「万」
 * - 小于 1 万亿：单位为「亿」
 * - 其他：单位为「万亿」
 * - 大于 1 万，且前缀数字小于 100 时，小数点后保留一位，为「0」时不展示
 * @param num
 */

func CnNiceNumber(num int64) string {
	if num < 10000 {
		return fmt.Sprintf("%d ", num) // 纯数字无文字时最后添加一个空格
	}
	if num < 1000000 {
		return fmt.Sprintf("%s 万", removePointZero(toFixedOneRound(num, 10000)))
	}
	if num < 100000000 {
		return fmt.Sprintf("%.0f 万", math.Floor(float64(num)/10000))
	}
	if num < 10000000000 {
		return fmt.Sprintf("%s 亿", removePointZero(toFixedOneRound(num, 100000000)))
	}
	if num < 1000000000000 {
		return fmt.Sprintf("%.0f 亿", math.Floor(float64(num)/100000000))

	}
	if num < 100000000000000 {
		return fmt.Sprintf("%s 万亿", removePointZero(toFixedOneRound(num, 1000000000000)))
	}
	return fmt.Sprintf("%.0f 万亿", math.Floor(float64(num)/1000000000000))
}

func toFixedOneRound(num, baseNum int64) string {
	r := fmt.Sprintf("%.1f", round(math.Floor((float64(num)/float64(baseNum))*10)/10, 0.1))
	return r
}

// For details, see https://stackoverflow.com/a/39544897/1705598
func round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func removePointZero(str string) string {
	index := strings.LastIndex(str, ".0")
	if index == len(str)-2 {
		return str[:index]
	}
	return str
}
