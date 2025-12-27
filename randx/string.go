package randx

import (
	"fmt"
	"strings"
)

// String 随机字符串
func String(length ...int) string {
	var l = 1 << 5
	if len(length) > 0 && length[0] > 0 {
		l = length[0]
	}
	return string(newRunes(l, []rune(allChar)))
}

// 生成随机字符
func newRunes(length int, pool []rune) []rune {
	runes, l := make([]rune, length), len(pool)
	for i := 0; i < length; i++ {
		index := IntRange(0, l-1)
		runes[i] = pool[index]
	}
	return runes
}

// SelectRune 选择字符
func SelectRune(pool []rune) rune {
	if l := len(pool); l > 0 {
		return pool[IntRange(0, l-1)]
	}
	return 0
}

// SelectString 选择字符串
func SelectString(pool []string) string {
	if l := len(pool); l > 0 {
		return pool[IntRange(0, l-1)]
	}
	return ""
}

type Option int

const (
	WithNumber        Option = 1 << 0 // 数字
	WithLowerLetter   Option = 1 << 1 // 小写字母
	WithUpperLetter   Option = 1 << 2 // 大写字母
	WithSpecialSymbol Option = 1 << 3 // 特殊符号
)

// StringWithOption 根据 Option 生成包含不同类型字符的字符串
func StringWithOption(length int, opt Option) string {
	var pool = numbers
	if opt&WithLowerLetter > 0 {
		pool += lowerLetters
	}
	if opt&WithUpperLetter > 0 {
		pool += upperLetters
	}
	if opt&WithSpecialSymbol > 0 {
		pool += special
	}

	runes := newRunes(length, []rune(pool))
	for i := range runes {
		j := NewRand().Intn(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// NumberCode 随机长度数字码
func NumberCode(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		y := IntRange(0, len(numbers)-1)
		bytes[i] = numbers[y]
	}
	return string(bytes)
}

// Name 随机姓名
func Name() string {
	sb := strings.Builder{}
	sb.WriteRune(SelectRune([]rune("赵钱孙李周吴郑王黄高冯陈蒋沈韩杨朱秦许何吕张孔曹")))
	sb.WriteRune(chineseRune())
	if Bool() {
		sb.WriteRune(chineseRune())
	}
	return sb.String()
}

// 随机汉字字符
func chineseRune() rune {
	return rune(IntRange(0x4E00, 0x9FA5))
}

// Chinese 随机汉字
func Chinese(length ...int) string {
	var l = 1
	if len(length) > 0 && length[0] > 0 {
		l = length[0]
	}
	chinese := make([]rune, l)
	for i := 0; i < l; i++ {
		chinese[i] = chineseRune()
	}
	return string(chinese)
}

// Phone 随机手机号
func Phone() string {
	bytes := make([]byte, 11)
	bytes[0] = '1'
	for i := 1; i < 11; i++ {
		y := IntRange(0, len(numbers)-1)
		bytes[i] = numbers[y]
	}
	return string(bytes)
}

// Email 随机邮箱号
func Email() string {
	sb := strings.Builder{}
	x, y := IntRange(5, 10), IntRange(2, 5)
	for i := 0; i < x; i++ {
		sb.WriteRune(SelectRune([]rune(lowerChar)))
	}
	sb.WriteString(`@`)
	for i := 0; i < y; i++ {
		sb.WriteRune(SelectRune([]rune(lowerLetters)))
	}
	sb.WriteString(`.com`)
	return sb.String()
}

func Ip() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		IntRange(1, 255),
		IntRange(0, 255),
		IntRange(0, 255),
		IntRange(0, 255))
}
