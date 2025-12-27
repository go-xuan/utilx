package marshalx

import (
	"fmt"
	"testing"
)

type demo struct {
	Id          int64  `properties:"id"`
	Name        string `properties:"name"`
	Default     string `properties:"default"`
	Description string `properties:"description"`
}

func TestMarshal(t *testing.T) {
	var d = demo{
		Id:          1123,
		Name:        "bbb",
		Default:     "aaa",
		Description: "abcd",
	}
	// json
	fmt.Print("json\n")
	data, _ := Json().Marshal(d)
	fmt.Println(string(data))

	// yaml
	fmt.Print("\nyaml\n")
	data, _ = Yaml().Marshal(d)
	fmt.Println(string(data))

	// msgpack
	fmt.Print("\nMsgpack\n")
	data, _ = Msgpack().Marshal(d)
	fmt.Println(string(data))

	fmt.Print("\nproperties\n")
	data, _ = Properties().Marshal(d)
	fmt.Println(string(data))

	fmt.Print("\ntoml\n")
	data, _ = Toml().Marshal(d)
	fmt.Println(string(data))
}
