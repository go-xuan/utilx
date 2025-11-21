package marshalx

import (
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
	"github.com/go-xuan/utilx/stringx"
)

const (
	JSON       = "json"
	YML        = "yml"
	YAML       = "yaml"
	TOML       = "toml"
	PROPERTIES = "properties"
	MSGPACK    = "msgpack"
)

// Marshal 序列化
type Marshal interface {
	Name() string
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	Read(string, interface{}) error
	Write(string, interface{}) error
}

// Apply 适配序列化
func Apply(name string) Marshal {
	if stringx.ContainsAny(name, ".", "\\", "/") {
		name = filex.GetSuffix(name)
	}
	switch name {
	case YML, YAML:
		return Yaml()
	case TOML:
		return Toml()
	case PROPERTIES:
		return Properties()
	case MSGPACK:
		return Msgpack()
	case JSON:
		return Json("    ")
	default:
		return Json()
	}
}

func readFile(path string) ([]byte, error) {
	if !filex.Exists(path) {
		return nil, errorx.Newf("the file not exist: %s", filex.Pwd(path))
	} else if data, err := filex.ReadFile(path); err != nil {
		return nil, errorx.Wrap(err, "read file error")
	} else {
		return data, nil
	}
}
