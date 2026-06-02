package stringx

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ParseInt 解析数字
func ParseInt(str string, def ...int) int {
	if value, err := strconv.Atoi(str); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// ParseInt64 解析数字
func ParseInt64(str string, def ...int64) int64 {
	if value, err := strconv.ParseInt(str, 10, 64); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// ParseUint64 解析无符号数字
func ParseUint64(str string, def ...uint64) uint64 {
	if value, err := strconv.ParseUint(str, 10, 64); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// ParseFloat 解析浮点数
func ParseFloat(str string, def ...float64) float64 {
	if value, err := strconv.ParseFloat(str, 64); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// ParseBool 解析布尔值
func ParseBool(str string, def ...bool) bool {
	str = strings.ToLower(str)
	for _, item := range []string{"1", "true", "是", "yes"} {
		if str == item {
			return true
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

// ParseTime 解析时间字符串
func ParseTime(str string, def ...time.Time) time.Time {
	if len(str) == 10 && str[4:5] == "-" {
		if parse, err := time.Parse("2006-01-02", str); err == nil {
			return parse
		}
	}
	if parse, err := time.Parse("2006-01-02 15:04:05", str); err == nil {
		return parse
	}
	if len(def) > 0 {
		return def[0]
	}
	return time.Time{}
}

// Int 转为字符串
func Int(i int) string {
	return strconv.Itoa(i)
}

// Int64 转为字符串
func Int64(i int64) string {
	return strconv.FormatInt(i, 10)
}

// Float 转为字符串
func Float(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// Json 转为 JSON 字符串
func Json(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

// JsonIndent 转为 JSON 字符串（格式化）
func JsonIndent(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "  ")
	return string(data)
}

// Yaml 转为 YAML 字符串
func Yaml(v interface{}) string {
	data, _ := yaml.Marshal(v)
	return string(data)
}

// Between 获取起始字符首次出现和结尾字符末次出现的下标
func Between(str, start, end string) (from, to int) {
	from, to = -1, -1
	if start == end {
		if indices := Indices(str, start, 2); len(indices) == 2 {
			from, to = indices[0], indices[1]
		} else if len(indices) == 1 {
			from = indices[0]
		}
		return
	}
	a, b, c := []rune(str), []rune(start), []rune(end)
	l, m, n := len(a), len(b), len(c)
	if m > l || n > l {
		return
	}
	var x, y int
	for i := 0; i < l; i++ {
		if a[i] == b[0] && string(a[i:i+m]) == start {
			if x++; x == 1 {
				from = i
			}
			i = i + m - 1
			continue
		}
		if x > 0 && a[i] == c[0] && string(a[i:i+n]) == end {
			if y++; y == x {
				to = i
				break
			}
			i = i + n - 1
		}
	}
	if to == -1 {
		from = -1
	}
	return
}

// Indices 获取所有下标，size：命中数量
func Indices(str, sub string, size ...int) []int {
	var limit int
	if len(size) > 0 {
		limit = size[0]
	}
	a, b := []rune(str), []rune(sub)
	l, m := len(a), len(b)
	indices := make([]int, 0)
	for i, n := 0, 0; i <= l-m; i++ {
		if limit > 0 && n >= limit {
			break
		}
		if a[i] == b[0] && string(a[i:i+m]) == sub {
			indices = append(indices, i)
			i = i + m - 1
			n++
		}
	}
	return indices
}

// Index 获取子串的下标
// position：表示获取位置，默认 position=1 即正序第 1 处，position=-1 即倒序第 1 处
func Index(str, sub string, position ...int) int {
	a, b := []rune(str), []rune(sub)
	l, m := len(a), len(b)
	switch {
	case m == 0:
		return 0
	case m == l && str == sub:
		return 0
	case m < l:
		x, y := 1, 0
		if len(position) > 0 {
			x = position[0]
		}
		for i := 0; i <= l-m; i++ {
			if x > 0 {
				if a[i] == b[0] && string(a[i:i+m]) == sub {
					if y++; x == y {
						return i
					}
				}
			} else {
				if a[l-i-1] == b[m-1] && string(a[l-i-m:l-i]) == sub {
					if y--; x == y {
						return l - i - m
					}
				}
			}
		}
	}
	return -1
}

// IndexStrict 获取子串下标（严格模式：仅当子串是独立单词时才命中）
func IndexStrict(str, key string) int {
	length, loop, index := len(key), true, 0
	for loop {
		if i := Index(str, key, 1); i >= 0 {
			if HasAdjacent(str, key, " ", i) {
				index, loop = index+i, false
			} else {
				index = i + length
				str = str[index:]
			}
		} else {
			index, loop = -1, false
		}
	}
	return index
}

// HasAdjacent 判断目标 key 在文本中当前位置是否有相邻字符
func HasAdjacent(str, key, adjacent string, index int) bool {
	sl, kl, al := len(str), len(key), len(adjacent)
	if index == 0 {
		return str[kl:kl+al] == adjacent
	} else if index == sl-kl {
		return str[index-al:index] == adjacent
	}
	return str[index-al:index] == adjacent && str[index+kl:index+kl+al] == adjacent
}

// AddPrefix 添加前缀
func AddPrefix(str, prefix string) string {
	if strings.HasPrefix(str, prefix) {
		return str
	}
	return prefix + str
}

// AddSuffix 添加后缀
func AddSuffix(str, suffix string) string {
	if strings.HasSuffix(str, suffix) {
		return str
	}
	return str + suffix
}

// Split 字符串分割（自动去除空白）
func Split(str string, sep string) []string {
	slice := strings.Split(str, sep)
	for i, s := range slice {
		slice[i] = strings.TrimSpace(s)
	}
	return slice
}

// Spaces 生成指定长度的空格字符串
func Spaces(length int) string {
	return strings.Repeat(" ", length)
}

// Contains 字符串是否包含子串
func Contains(str string, substr string) bool {
	return strings.Contains(str, substr)
}

// ContainsAny 字符串 str 是否包含 keys 中的任意值
func ContainsAny(str string, substr ...string) bool {
	for _, key := range substr {
		if strings.Contains(str, key) {
			return true
		}
	}
	return false
}

// ContainsAll 字符串 str 是否包含 keys 中的所有值
func ContainsAll(str string, substr ...string) bool {
	for _, sep := range substr {
		if !strings.Contains(str, sep) {
			return false
		}
	}
	return true
}

// HasEmpty 是否有空字符串
func HasEmpty(str ...string) bool {
	for _, item := range str {
		if item == "" {
			return true
		}
	}
	return false
}

// Default 用于函数中的不定参数取默认值
func Default(def string, variable ...string) string {
	if len(variable) > 0 && variable[0] != "" {
		return variable[0]
	}
	return def
}

// IfZero 为空时取默认值
func IfZero(value, def string) string {
	if value == "" {
		return def
	}
	return value
}

// Reverse 反转字符串
func Reverse(str string) string {
	runes := []rune(str)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// SubString 字符串截取
func SubString(str string, start, end int) string {
	r := []rune(str)
	length := len(r)
	if end > length {
		end = length
	}
	if length > 0 && start >= 0 && start <= end {
		return string(r[start:end])
	}
	return ""
}

// Cut 分割字符串
// position：表示分割位置，默认 position=1 即正序第 1 处，position=-1 即倒序第 1 处
func Cut(str, cut string, position ...int) (string, string) {
	if i := Index(str, cut, position...); i >= 0 {
		return str[:i], str[i+len(cut):]
	}
	return str, ""
}

// Insert 插入字符串
func Insert(str, insert string, position ...int) string {
	if len(position) > 0 {
		if i := position[0]; i > 0 && i < len(str) {
			return str[:i] + insert + str[i:]
		}
	}
	return str + insert
}

// Fill 字符填充
func Fill(str, fill string, length int) string {
	l := len(str)
	if length > 0 {
		return str + Grow(fill, length-l)
	}
	return Grow(fill, -length-l) + str
}

// Grow 字符扩充到固定长度
func Grow(str string, length int) string {
	l := len(str)
	if length <= 0 || l == 0 {
		return str
	} else if l == 1 {
		return strings.Repeat(str, length)
	}
	var sb strings.Builder
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(str[i%l])
	}
	return sb.String()
}

// Reduce 字符串缩减（重复子串）
func Reduce(str string) (string, bool) {
	if l := len(str); l == 2 {
		if str[0] == str[1] {
			return str[:1], true
		}
	} else if l > 2 {
		for i := 1; i <= l/2; i++ {
			if l%i == 0 {
				sub := str[:i]
				if strings.Repeat(sub, l/i) == str {
					return sub, true
				}
			}
		}
	}
	return str, false
}

// ParseUrlParams 解析 url 参数为 map
func ParseUrlParams(str string) map[string]string {
	if str == "" {
		return nil
	}
	params := make(map[string]string)
	for _, kv := range strings.Split(str, "&") {
		k, v := Cut(kv, "=")
		params[k] = v
	}
	return params
}

// ToSnake 转下划线
func ToSnake(str string) string {
	runes, j := []rune(str), false
	var res []rune
	for i, r := range runes {
		if i > 0 && r >= 65 && r <= 90 && j {
			res = append(res, 95)
		} else if r != 95 {
			j = true
		}
		res = append(res, r)
	}
	return strings.ToLower(string(res))
}

// ToLowerCamel 转小驼峰
func ToLowerCamel(str string) string {
	runes := toUpperCamel([]rune(str))
	if len(runes) > 0 {
		runes[0] = runes[0] + 32
	}
	return string(runes)
}

// ToUpperCamel 转大驼峰
func ToUpperCamel(str string) string {
	runes := toUpperCamel([]rune(str))
	return string(runes)
}

func toUpperCamel(runes []rune) []rune {
	up := true
	var res []rune
	for _, r := range runes {
		if r == 95 {
			up = true
			continue
		} else if r >= 65 && r <= 90 && up {
			up = false
		} else if r >= 65 && r <= 90 {
			r = r + 32
		} else if r >= 97 && r <= 122 && up {
			r, up = r-32, false
		}
		res = append(res, r)
	}
	return res
}

// Similarity 文本相似度计算（Levenshtein 距离）
func Similarity(source, target string) float64 {
	sl, tl := len(source), len(target)
	if (sl == 0 && tl == 0) || source == target {
		return 1.0
	}
	matrix := make([][]int, sl+1)
	for i := range matrix {
		matrix[i] = make([]int, tl+1)
		matrix[i][0] = i
	}

	for j := 0; j <= tl; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= sl; i++ {
		for j := 1; j <= tl; j++ {
			cost := 0
			if source[i-1] != target[j-1] {
				cost = 1
			}
			matrix[i][j] = minOfThree(matrix[i-1][j]+1, matrix[i][j-1]+1, matrix[i-1][j-1]+cost)
		}
	}

	distance := matrix[sl][tl]
	maxLen := float64(sl)
	if tl > sl {
		maxLen = float64(tl)
	}
	return 1.0 - float64(distance)/maxLen
}

func minOfThree(a, b, c int) int {
	if a <= b && a <= c {
		return a
	} else if b <= a && b <= c {
		return b
	}
	return c
}

// MatchUrl URL 匹配
func MatchUrl(uri, rule string) bool {
	if rule == "*" || rule == "/*" {
		return true
	} else if strings.HasSuffix(rule, `/**`) {
		return strings.HasPrefix(uri, rule[:len(rule)-3])
	} else if strings.HasSuffix(rule, `/*`) {
		prefix := rule[:len(rule)-2]
		if strings.HasPrefix(uri, prefix) {
			uri = uri[len(prefix):]
			return Index(uri, `/`) < 0
		}
	} else {
		return strings.Contains(uri, rule)
	}
	return false
}
