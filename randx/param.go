package randx

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-xuan/utilx/stringx"
	"github.com/go-xuan/utilx/timex"
	"github.com/google/uuid"
)

func NewParam(data map[string]string) *Param {
	param := &Param{}
	if data != nil {
		param.Min = data["min"]
		param.Max = data["max"]
		param.Prefix = data["prefix"]
		param.Suffix = data["suffix"]
		param.Upper = data["upper"] == "true"
		param.Lower = data["lower"] == "true"
		param.Old = data["old"]
		param.New = data["new"]
		param.Length = stringx.ParseInt(data["length"])
		param.Prec = stringx.ParseInt(data["prec"])
		param.Level = stringx.ParseInt(data["level"])
		param.Enums = strings.Split(data["enums"], ",")
	}
	return param
}

// Param 随机数生成参数
type Param struct {
	Min    string   // 最小值
	Max    string   // 最大值
	Prefix string   // 前缀
	Suffix string   // 后缀
	Upper  bool     // 转大写
	Lower  bool     // 转小写
	Old    string   // 替换旧字符
	New    string   // 替换新字符
	Length int      // 长度
	Prec   int      // 小数位精度
	Level  int      // 级别
	Enums  []string // 枚举选项，多个以逗号分割
}

// 字符串
func (p *Param) String(type_ string) string {
	switch type_ {
	case INT:
		return strconv.Itoa(p.Int())
	case FLOAT:
		return strconv.FormatFloat(p.Float(), 'f', -1, 64)
	case TIME:
		return timex.Format(p.Time())
	case DATE:
		return timex.Format(p.Time(), timex.DateFmt)
	case UUID:
		return uuid.NewString()
	case PHONE:
		return Phone()
	case NAME:
		return Name()
	case EMAIL:
		return Email()
	case IP:
		return Ip()
	case PASSWORD:
		return p.Password()
	case ENUM:
		return SelectString(p.Enums)
	default:
		return String(p.Length)
	}
}

// Int 整数
func (p *Param) Int() int {
	min_ := stringx.ParseInt(p.Min, 1)
	max_ := stringx.ParseInt(p.Max, 999)
	return IntRange(min_, max_)
}

// Float 浮点数
func (p *Param) Float() float64 {
	min_ := stringx.ParseFloat(p.Min, 1)
	max_ := stringx.ParseFloat(p.Max, 999)
	if p.Prec == 0 {
		p.Prec = 6
	}
	return Float64Range(min_, max_, p.Prec)
}

// Time 时间
func (p *Param) Time() time.Time {
	start := stringx.ParseTime(p.Min, time.Time{})
	end := stringx.ParseTime(p.Max, time.Now())
	return TimeRange(start, end)
}

// Enum 枚举值
func (p *Param) Enum() string {
	return SelectString(p.Enums)
}

// Sequence 序列值
func (p *Param) Sequence(offset int) int {
	return stringx.ParseInt(p.Min, 1) + offset
}

// Password 密码
func (p *Param) Password() string {
	switch p.Level {
	case 2:
		return StringWithOption(p.Length, WithNumber|WithLowerLetter|WithUpperLetter)
	case 3:
		return StringWithOption(p.Length, WithNumber|WithLowerLetter|WithUpperLetter|WithSpecialSymbol)
	default:
		return StringWithOption(p.Length, WithNumber)
	}
}

// AddEnum 添加枚举项
func (p *Param) AddEnum(enum ...string) {
	p.Enums = append(p.Enums, enum...)
}

// Modify 修饰
func (p *Param) Modify(value string) string {
	if p.Old != "" {
		value = strings.ReplaceAll(value, p.Old, p.New)
	}
	if p.Prefix != "" {
		value = p.Prefix + value
	}
	if p.Suffix != "" {
		value = value + p.Suffix
	}
	if p.Upper {
		value = strings.ToUpper(value)
	}
	if p.Lower {
		value = strings.ToLower(value)
	}
	return value
}
