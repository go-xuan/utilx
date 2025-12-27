package marshalx

import (
	"gopkg.in/yaml.v3"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
)

func Yaml() Marshal {
	return yamlImpl{}
}

type yamlImpl struct{}

func (y yamlImpl) Name() string {
	return YAML
}

func (y yamlImpl) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y yamlImpl) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (y yamlImpl) Read(path string, v interface{}) error {
	data, err := readFile(path)
	if err != nil {
		return errorx.Wrap(err, "read file error")
	}
	return y.Unmarshal(data, v)
}

func (y yamlImpl) Write(path string, v interface{}) error {
	if data, err := y.Marshal(v); err != nil {
		return errorx.Wrap(err, "yaml marshal error")
	} else if err = filex.WriteFile(path, data); err != nil {
		return errorx.Wrap(err, "write yaml file error")
	}
	return nil
}
