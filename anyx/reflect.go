package anyx

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/stringx"
)

// ValueOf 反射获取结构体的值
func ValueOf(v any) reflect.Value {
	var valueOf = reflect.ValueOf(v)
	if valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	return valueOf
}

// TypeOf 反射获取结构体的类型
func TypeOf(v any) reflect.Type {
	var typeOf = reflect.TypeOf(v)
	if typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
	}
	return typeOf
}

// GetTagValues 获取结构体中指定tag的值
func GetTagValues(v any, tag string) []string {
	valueOf := ValueOf(v)
	typeOf := valueOf.Type()

	var names []string
	for i := 0; i < typeOf.NumField(); i++ {
		if name, ok := typeOf.Field(i).Tag.Lookup(tag); ok && name != "" {
			names = append(names, name)
		}
	}
	return names
}

// GetFieldValuesWithTag 获取结构体中包含指定tag的字段值
func GetFieldValuesWithTag(v any, tag string) []string {
	valueOf := ValueOf(v)
	typeOf := valueOf.Type()

	if typeOf.Kind() != reflect.Struct {
		return nil
	}
	var values []string
	for i := 0; i < typeOf.NumField(); i++ {
		if name, ok := typeOf.Field(i).Tag.Lookup(tag); ok && name != "" {
			values = append(values, fmt.Sprintf("%v", valueOf.Field(i).Interface()))
		}
	}
	return values
}

// SetFieldValueFromTag 将结构体的tag值设为字段值
func SetFieldValueFromTag(v any, tag string) error {
	if err := MustStructPointer(v); err != nil {
		return errorx.New("the kind must be struct pointer")
	}
	valueOf := ValueOf(v)
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		if field.IsZero() {
			if value := valueOf.Type().Field(i).Tag.Get(tag); value != "" {
				if err := SetValue(field, value); err != nil {
					return errorx.Wrap(err, "failed to set default value")
				}
			}
		}
	}
	return nil
}

// MustStructPointer 检查是否为结构体指针
func MustStructPointer(v any) error {
	var typeOf = reflect.TypeOf(v)
	if typeOf.Kind() != reflect.Pointer {
		return errorx.New("the kind must be pointer")
	} else if typeOf.Elem().Kind() == reflect.Slice {
		return errorx.New("the kind cannot be slice")
	}
	return nil
}

// SetValue 反射设置字段值
func SetValue(value reflect.Value, str string) error {
	switch k := value.Kind(); k {
	case reflect.String:
		value.SetString(str)
	case reflect.Bool:
		value.SetBool(stringx.ParseBool(str))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, _ := strconv.ParseInt(str, 10, 64)
		value.SetInt(x)
	case reflect.Float32, reflect.Float64:
		float, _ := strconv.ParseFloat(str, 64)
		value.SetFloat(float)
	default:
		return errorx.New(fmt.Sprintf("unsupported kind %T", k.String()))
	}
	return nil
}

// MapToStruct 将map转换为结构体
func MapToStruct(m map[string]string, v any) error {
	valueOf := ValueOf(v)
	for key, value := range m {
		field := valueOf.FieldByName(key)
		if field.IsValid() && field.CanSet() {
			if err := SetValue(field, value); err != nil {
				return errorx.Wrap(err, fmt.Sprintf("set value to field %s error", key))
			}
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
		fieldA := va.Type().Field(i)
		fieldB := vb.FieldByName(fieldA.Name)
		if fieldB.IsValid() && fieldB.CanSet() {
			va.Field(i).Set(fieldB)
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
