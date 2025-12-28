package marshalx

import (
	"fmt"
	"testing"

	"github.com/go-xuan/utilx/stringx"
)

type Demo struct {
	Id   int64  `properties:"id"`
	Name string `properties:"name"`
	Sub1 Sub    `properties:"sub1"`
	Sub2 *Sub   `properties:"sub2"`
	Sub3 []Sub  `properties:"sub3"`
	Sub4 []*Sub `properties:"sub4"`
}

type Sub struct {
	Id   int64  `properties:"id"`
	Name string `properties:"name"`
}

func newDemo() Demo {
	return Demo{
		Id:   0,
		Name: "demo",
		Sub1: Sub{Id: 1, Name: "Sub1"},
		Sub2: &Sub{Id: 2, Name: "Sub2"},
		Sub3: []Sub{{Id: 30, Name: "Sub3-0"}, {Id: 31, Name: "Sub3-1"}},
		Sub4: []*Sub{{Id: 40, Name: "Sub4-0"}, {Id: 41, Name: "Sub4-1"}},
	}
}

func TestMarshal(t *testing.T) {
	demo := newDemo()
	t.Run("marshal json", func(t *testing.T) {
		data, _ := Json().Marshal(demo)
		fmt.Println(string(data))
	})
	t.Run("marshal yaml", func(t *testing.T) {
		data, _ := Yaml().Marshal(demo)
		fmt.Println(string(data))
	})
	t.Run("marshal msgpack", func(t *testing.T) {
		data, _ := Msgpack().Marshal(demo)
		fmt.Println(string(data))
	})
	t.Run("marshal properties", func(t *testing.T) {
		data, _ := Properties().Marshal(demo)
		fmt.Println(string(data))
	})
	t.Run("marshal toml", func(t *testing.T) {
		data, _ := Toml().Marshal(demo)
		fmt.Println(string(data))
	})
}

func TestUnmarshal(t *testing.T) {
	demo := Demo{}
	t.Run("marshal properties", func(t *testing.T) {
		data := []byte(`
id=0
name=demo
sub1.id=1
sub1.name=Sub1
sub2.id=2
sub2.name=Sub2
sub3[0].id=30
sub3[0].name=Sub3-0
sub3[1].id=31
sub3[1].name=Sub3-1
sub4[0].id=40
sub4[0].name=Sub4-0
sub4[1].id=41
sub4[1].name=Sub4-1
`)
		_ = Properties().Unmarshal(data, &demo)
		fmt.Println(stringx.JsonIndent(demo))
	})
}
