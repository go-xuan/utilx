package marshalx

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/magiconair/properties"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
	"github.com/go-xuan/utilx/reflectx"
)

func Properties() Marshal {
	return propertiesImpl{}
}

type propertiesImpl struct{}

func (p propertiesImpl) Name() string {
	return PROPERTIES
}

func (p propertiesImpl) Marshal(v interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	propertiesMarshal(buffer, reflectx.ValueOf(v), "")
	return buffer.Bytes(), nil
}

func (p propertiesImpl) Unmarshal(data []byte, v interface{}) error {
	if !reflectx.IsStructPointer(v) {
		return errorx.New("the kind must be struct pointer")
	}
	pp, err := properties.Load(data, properties.UTF8)
	if err != nil {
		return errorx.Wrap(err, "load properties error")
	}
	propertiesSetStruct(pp, reflectx.ValueOf(v), "")
	return nil
}

// Read 读取properties文件
func (p propertiesImpl) Read(path string, v interface{}) error {
	data, err := readFile(path)
	if err != nil {
		return errorx.Wrap(err, "read file error")
	}
	return p.Unmarshal(data, v)
}

// Write 写入properties文件
func (p propertiesImpl) Write(path string, v interface{}) error {
	if data, err := p.Marshal(v); err != nil {
		return errorx.Wrap(err, "properties marshal error")
	} else if err = filex.WriteFile(path, data); err != nil {
		return errorx.Wrap(err, "write properties file error")
	}
	return nil
}

// properties序列化
func propertiesMarshal(b *bytes.Buffer, value reflect.Value, parent string) {
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		fieldValue, structField := value.Field(i), typ.Field(i)
		if key := structField.Tag.Get(PROPERTIES); key != "" {
			key = propertiesKey(key, parent)
			switch fieldValue.Kind() {
			case reflect.Struct, reflect.Pointer:
				propertiesMarshal(b, fieldValue, key)
			case reflect.Slice:
				propertiesMarshalSlice(b, fieldValue, key)
			default:
				b.WriteString(fmt.Sprintf("%s=%v\n", key, fieldValue.Interface()))
			}
		}
	}
}

// properties序列化slice
func propertiesMarshalSlice(b *bytes.Buffer, value reflect.Value, parent string) {
	for i := 0; i < value.Len(); i++ {
		propertiesMarshal(b, value.Index(i), fmt.Sprintf("%s[%d]", parent, i))
	}
}

// properties赋值给结构体
func propertiesSetStruct(p *properties.Properties, value reflect.Value, parent string) {
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		fieldValue, structField := value.Field(i), typ.Field(i)
		if key := structField.Tag.Get(PROPERTIES); key != "" && fieldValue.CanSet() {
			key = propertiesKey(key, parent)
			switch fieldValue.Kind() {
			case reflect.String:
				fieldValue.SetString(p.GetString(key, ""))
			case reflect.Bool:
				fieldValue.SetBool(p.GetBool(key, false))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue.SetInt(int64(p.GetInt(key, 0)))
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue.SetUint(uint64(p.GetInt(key, 0)))
			case reflect.Float32, reflect.Float64:
				fieldValue.SetFloat(p.GetFloat64(key, 0))
			case reflect.Struct:
				propertiesSetStruct(p, fieldValue, key)
			case reflect.Pointer:
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(structField.Type.Elem()))
				}
				propertiesSetStruct(p, fieldValue.Elem(), key)
			case reflect.Slice, reflect.Array:
				propertiesSetSlice(p, fieldValue, key)
			default:
			}
		}
	}
}

// properties赋值slice
func propertiesSetSlice(p *properties.Properties, value reflect.Value, parent string) {
	elem, pointer := value.Type().Elem(), false
	if elem.Kind() == reflect.Pointer {
		elem, pointer = elem.Elem(), true
	}
	index := 0
	for {
		key := fmt.Sprintf("%s[%d]", parent, index)
		if !propertiesExist(p, elem, key) {
			break
		}
		instance := reflect.New(elem).Elem()
		propertiesSetStruct(p, instance, key)
		if pointer {
			instance = instance.Addr()
		}
		value.Set(reflect.Append(value, instance))
		index++
	}
}

// 获取properties key
func propertiesKey(key string, parent string) string {
	if parent != "" {
		key = fmt.Sprintf("%s.%s", parent, key)
	}
	return key
}

// 检查是否存在properties键值对
func propertiesExist(p *properties.Properties, typ reflect.Type, parent string) bool {
	for i := 0; i < typ.NumField(); i++ {
		if key := typ.Field(i).Tag.Get(PROPERTIES); key != "" {
			key = propertiesKey(key, parent)
			if _, ok := p.Get(key); ok {
				return true
			}
		}
	}
	return false
}
