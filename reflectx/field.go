package reflectx

import (
	"reflect"

	"github.com/go-xuan/typex"
)

// Fields 获取结构体字段
func Fields(v any) []*Field {
	if IsStruct(v) {
		val := ValueOf(v)
		typ := val.Type()
		var fields []*Field
		for i := 0; i < typ.NumField(); i++ {
			fields = append(fields, &Field{
				Info:  typ.Field(i),
				Value: val.Field(i),
			})
		}
		return fields
	}
	return nil
}

// ExtractValues 提取结构体字段的value
func ExtractValues(v any, tag ...string) ([]string, []string) {
	if fields := Fields(v); len(fields) > 0 {
		var names, values []string
		for _, field := range fields {
			if name, value := field.Extract(tag...); name != "" {
				names = append(names, name)
				values = append(values, value)
			}
		}
		return names, values
	}
	return nil, nil
}

// Rows 获取结构体字段的name-value对
func Rows[T any](slice []T, withHeader bool) [][]string {
	var rows [][]string
	for i, item := range slice {
		names, values := ExtractValues(item)
		if i == 0 && withHeader {
			rows = append(rows, names)
		}
		rows = append(rows, values)
	}
	return rows
}

// Field 结构体字段
type Field struct {
	Info  reflect.StructField // 字段信息
	Value reflect.Value       // 反射值
}

// GetName 获取结构体字段的名称
func (f *Field) GetName() string {
	return f.Info.Name
}

// GetValue 获取结构体字段的值
func (f *Field) GetValue() typex.Value {
	return GetValue(f.Value)
}

// GetKind 获取结构体字段的类型
func (f *Field) GetKind() reflect.Kind {
	return f.Info.Type.Kind()
}

// GetInterface 获取结构体字段的值
func (f *Field) GetInterface() interface{} {
	if f.Info.IsExported() {
		return f.Value.Interface()
	}
	return nil
}

// TagLookup 获取结构体字段的标签
func (f *Field) TagLookup(tag string) (string, bool) {
	return f.Info.Tag.Lookup(tag)
}

// Extract 提取name-value对
func (f *Field) Extract(tag ...string) (string, string) {
	if len(tag) == 0 || tag[0] == "" {
		return f.GetName(), f.GetValue().String()
	} else if name, _ := f.TagLookup(tag[0]); name != "-" {
		return name, f.GetValue().String()
	}
	return "", ""
}
