package stringx

import (
	"unicode"
)

// Align 对齐字符串
func Align(str string, length int) string {
	switch l := len(str); {
	case length > 0 && length > l:
		return str + Spaces(length-l)
	case length < 0 && -length > l:
		return Spaces(-length-l) + str
	default:
		return str
	}
}

// VisualLength 计算字符串的可视化长度（1个中文占5/3个字符宽度）
func VisualLength(s string) int {
	var runes, cn = []rune(s), 0
	for _, r := range runes {
		if unicode.Is(unicode.Han, r) {
			cn++
		}
	}
	length := len(runes) + cn*2/3
	if y := cn % 3; y == 1 {
		length++
	}
	return length
}

// AlignSlice 对齐字符串数组
func AlignSlice(slice []string, length int) {
	for i, value := range slice {
		slice[i] = Align(value, length)
	}
}

// AlignRows 对齐二维字符串数组
func AlignRows(rows [][]string) {
	lengths := MaxLengths(rows)
	for _, row := range rows {
		for i, value := range row {
			row[i] = Align(value, lengths[i])
		}
	}
}

// MaxLengths 计算二维字符串数组没列的最大可视化长度
func MaxLengths(rows [][]string) []int {
	var lengths []int
	for _, row := range rows {
		lengths = mergeMaxLengths(lengths, row)
	}
	return lengths
}

// 合并每个位置的最大集
func mergeMaxLengths(lengths []int, slice []string) []int {
	if len(lengths) == 0 {
		lengths = make([]int, len(slice))
	}
	for i, s := range slice {
		if l := VisualLength(s); lengths[i] < l {
			lengths[i] = l
		}
	}
	return lengths
}
