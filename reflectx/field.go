package reflectx

import (
	"reflect"

	"github.com/go-xuan/typex"
)

// GetRowsByNames 根据字段名获取切片数据
func GetRowsByNames[T any](slice []T, names ...string) [][]string {
	var rows [][]string
	for i, item := range slice {
		headers, values := GetValuesByNames(item, names...)
		if i == 0 {
			rows = append(rows, headers)
		}
		rows = append(rows, values)
	}
	return rows
}

// GetValuesByNames 根据字段名获取结构体的name+value
func GetValuesByNames(v any, names ...string) ([]string, []string) {
	if fields := GetFieldsByNames(v, names...); len(fields) > 0 {
		var headers, values []string
		for _, field := range fields {
			name, value := field.Extract()
			headers = append(headers, name)
			values = append(values, value)
		}
		return headers, values
	}
	return nil, nil
}

// GetFieldsByNames 根据字段名提取字段
func GetFieldsByNames(v any, names ...string) []*Field {
	if IsStruct(v) {
		val := ValueOf(v)
		typ := val.Type()
		var fields []*Field
		if len(names) == 0 {
			for i := 0; i < typ.NumField(); i++ {
				fields = append(fields, &Field{
					Info:  typ.Field(i),
					Value: val.Field(i),
				})
			}
			return fields
		}

		set := make(map[string]*Field)
		for _, name := range names {
			field := &Field{}
			fields = append(fields, field)
			set[name] = field
		}
		for i := 0; i < typ.NumField(); i++ {
			info := typ.Field(i)
			if field, ok := set[info.Name]; ok {
				field.Info = info
				field.Value = val.Field(i)
			}
		}
		return fields
	}
	return nil
}

// GetRowsByTag 根据tag获取切片数据
func GetRowsByTag[T any](slice []T, tag ...string) [][]string {
	var rows [][]string
	for i, item := range slice {
		headers, values := GetValuesByTag(item, tag...)
		if i == 0 {
			rows = append(rows, headers)
		}
		rows = append(rows, values)
	}
	return rows
}

// GetValuesByTag 根据tag提取结构体的name+value
func GetValuesByTag(v any, tag ...string) ([]string, []string) {
	if fields := GetFieldsByTag(v, tag...); len(fields) > 0 {
		var headers, values []string
		for _, field := range fields {
			name, value := field.Extract(tag...)
			headers = append(headers, name)
			values = append(values, value)
		}
		return headers, values
	}
	return nil, nil
}

// GetFieldsByTag 根据tag提取字段
func GetFieldsByTag(v any, tag ...string) []*Field {
	if IsStruct(v) {
		val := ValueOf(v)
		typ := val.Type()
		var fields []*Field
		if len(tag) == 0 || tag[0] == "" {
			for i := 0; i < typ.NumField(); i++ {
				fields = append(fields, &Field{
					Info:  typ.Field(i),
					Value: val.Field(i),
				})
			}
			return fields
		}
		key := tag[0]
		for i := 0; i < typ.NumField(); i++ {
			info := typ.Field(i)
			if value, ok := info.Tag.Lookup(key); ok && value != "-" {
				fields = append(fields, &Field{
					Info:  info,
					Value: val.Field(i),
				})
			}
		}
		return fields
	}
	return nil
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
	value := f.GetValue().String()
	if len(tag) > 0 && tag[0] != "" {
		if name, ok := f.TagLookup(tag[0]); ok && name != "-" {
			return name, value
		}
	}
	return f.GetName(), value
}
