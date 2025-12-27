package marshalx

import (
	"bytes"

	"github.com/BurntSushi/toml"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
)

func Toml() Marshal {
	return tomlImpl{}
}

type tomlImpl struct{}

func (t tomlImpl) Name() string {
	return TOML
}

func (t tomlImpl) Marshal(v interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	if err := toml.NewEncoder(&buffer).Encode(v); err != nil {
		return nil, errorx.Wrap(err, "encode toml failed")
	}
	return buffer.Bytes(), nil
}

func (t tomlImpl) Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

func (t tomlImpl) Read(path string, v interface{}) error {
	data, err := readFile(path)
	if err != nil {
		return errorx.Wrap(err, "read file error")
	}
	return t.Unmarshal(data, v)
}

func (t tomlImpl) Write(path string, v interface{}) error {
	if data, err := t.Marshal(v); err != nil {
		return errorx.Wrap(err, "toml marshal error")
	} else if err = filex.WriteFile(path, data); err != nil {
		return errorx.Wrap(err, "write toml file error")
	}
	return nil
}
