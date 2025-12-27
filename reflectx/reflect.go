package reflectx

import (
	"fmt"
	"reflect"
	"time"

	"github.com/go-xuan/typex"
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/stringx"
)

const (
	String  = "string"
	Int     = "int"
	Int8    = "int8"
	Int16   = "int16"
	Int32   = "int32"
	Int64   = "int64"
	Uint    = "uint"
	Uint8   = "uint8"
	Uint16  = "uint16"
	Uint32  = "uint32"
	Uint64  = "uint64"
	Float32 = "float32"
	Float64 = "float64"
	Bool    = "bool"
	Time    = "time.Time"
)

// ValueOf 反射获取结构体的值
func ValueOf(v any) reflect.Value {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	return val
}

// TypeOf 反射获取结构体的类型
func TypeOf(v any) reflect.Type {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ
}

// IsStructPointer 检查是否为结构体指针
func IsStructPointer(v any) bool {
	if t := reflect.TypeOf(v); t.Kind() != reflect.Pointer {
		return false
	} else if t.Elem().Kind() == reflect.Slice {
		return false
	}
	return true
}

// IsStruct 检查是否为结构体或结构体指针
func IsStruct(v any) bool {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Struct {
		return true
	} else if t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct {
		return true
	}
	return false
}

// SetValue 反射设置字段值
func SetValue(value reflect.Value, str string) {
	switch kind := value.Type().String(); kind {
	case String:
		value.SetString(str)
	case Bool:
		value.SetBool(stringx.ParseBool(str))
	case Int, Int8, Int16, Int32, Int64:
		value.SetInt(stringx.ParseInt64(str))
	case Uint, Uint8, Uint16, Uint32, Uint64:
		value.SetUint(stringx.ParseUint64(str))
	case Float32, Float64:
		value.SetFloat(stringx.ParseFloat(str))
	case Time:
		value.Set(reflect.ValueOf(stringx.ParseTime(str)))
	default:
		value.SetZero()
	}
}

// GetValue 反射获取字段值
func GetValue(value reflect.Value) typex.Value {
	switch kind := value.Type().String(); kind {
	case String:
		return typex.NewString(value.String())
	case Bool:
		return typex.NewBool(value.Bool())
	case Int, Int8, Int16, Int32, Int64:
		return typex.NewInt64(value.Int())
	case Uint, Uint8, Uint16, Uint32, Uint64:
		return typex.NewInt64(int64(value.Uint()))
	case Float32, Float64:
		return typex.NewFloat64(value.Float())
	case Time:
		return typex.NewTime(value.Interface().(time.Time))
	default:
		return typex.NewZero()
	}
}

// MapToStruct 将map转换为结构体
func MapToStruct(m map[string]string, v any) error {
	val := ValueOf(v)
	for key, value := range m {
		if field := val.FieldByName(key); field.IsValid() && field.CanSet() {
			SetValue(field, value)
		} else if !field.IsValid() {
			// 如果没有找到对应的字段则返回错误
			return errorx.New(fmt.Sprintf("no such field %s in the structure", key))
		}
	}
	return nil
}

// MergeStructs 合并结构体
func MergeStructs(a, b any) {
	va, vb := ValueOf(a), ValueOf(b)
	for i := 0; i < va.NumField(); i++ {
		name := va.Type().Field(i).Name
		field := vb.FieldByName(name)
		if field.IsValid() && field.CanSet() {
			va.Field(i).Set(field)
		}
	}
}

// SetZeroValue 设置结构体零值
func SetZeroValue[T any](a, b T) {
	va, vb := ValueOf(a), ValueOf(b)
	for i := 0; i < va.NumField(); i++ {
		if va.Field(i).IsZero() {
			va.Field(i).Set(vb.Field(i))
		}
	}
}
