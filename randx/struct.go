package randx

import (
	"reflect"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/reflectx"
)

// Struct 随机填充结构体字段值
func Struct(v any) error {
	if !reflectx.IsStructPointer(v) {
		return errorx.New("the kind must be struct pointer")
	}
	val := reflectx.ValueOf(v)
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		if field := val.Field(i); field.IsZero() {
			structField := typ.Field(i)
			if value := structField.Tag.Get("default"); value != "" {
				reflectx.SetValue(field, value)
			} else if type_ := structField.Tag.Get("type"); type_ != "" {
				value = NewParam(nil).String(type_)
				reflectx.SetValue(field, value)
			} else {
				SetValue(field)
			}
		}
	}
	return nil
}

// SetValue 设置随机值
func SetValue(value reflect.Value) {
	switch kind := value.Type().String(); kind {
	case reflectx.String:
		value.SetString(String())
	case reflectx.Bool:
		value.SetBool(Bool())
	case reflectx.Int, reflectx.Int8, reflectx.Int16, reflectx.Int32, reflectx.Int64:
		value.SetInt(Int64())
	case reflectx.Uint, reflectx.Uint8, reflectx.Uint16, reflectx.Uint32, reflectx.Uint64:
		value.SetUint(uint64(Int64()))
	case reflectx.Float32, reflectx.Float64:
		value.SetFloat(Float64())
	case reflectx.Time:
		value.Set(reflect.ValueOf(Time()))
	default:
		value.SetZero()
	}
}
