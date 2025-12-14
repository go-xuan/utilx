package randx

import (
	"strconv"
)

const (
	STRING   = "string"   // 字符串
	INT      = "int"      // 数字
	FLOAT    = "float"    // 浮点数
	SEQUENCE = "sequence" // 序列
	TIME     = "time"     // 时间
	DATE     = "date"     // 日期
	UUID     = "uuid"     // uuid
	PHONE    = "phone"    // 手机号
	NAME     = "name"     // 姓名
	EMAIL    = "email"    // 邮箱
	IP       = "ip"       // ip地址
	PASSWORD = "password" // 密码
	ENUM     = "enum"     // 枚举
)

// NewString 生成随机数
func NewString(opt *Options) string {
	var value string
	if opt.Type == SEQUENCE {
		value = strconv.Itoa(opt.Param.Sequence(opt.Offset))
	} else {
		value = opt.Param.String(opt.Type)
	}
	return opt.Param.Modify(value)
}

// Options 随机生成
type Options struct {
	Param   *Param // 约束条件参数
	Type    string // 数据类型
	Default string // 默认值
	Offset  int    // 偏移量
}

// NewString 生成随机数
func (o *Options) NewString() string {
	var value string
	if o.Type == SEQUENCE {
		value = strconv.Itoa(o.Param.Sequence(o.Offset))
	} else {
		value = o.Param.String(o.Type)
	}
	return o.Param.Modify(value)
}
