package anyx

import (
	"fmt"
	"testing"
)

type demo struct {
	Name  string  `json:"name" default:"abc"`
	Num   string  `json:"num" default:"123"`
	Float float64 `json:"float" default:"123.45"`
	Null  string  `json:"null"`
}

func TestSetValueByTag(t *testing.T) {
	var d = demo{}

	if err := SetFieldValueFromTag(&d, "default"); err != nil {
		fmt.Println(err)
	}
	fmt.Println("name = ", d.Name)
	fmt.Println("num = ", d.Num)
	fmt.Println("float = ", d.Float)
	fmt.Println("null = ", d.Null)
}

func TestGetValuesByTag(t *testing.T) {
	var d = demo{}

	//_ = SetFieldValueFromTag(&d, "default")

	values := GetFieldValuesWithTag(d, "default")
	fmt.Println(values)

	names := GetTagValues(d, "default")
	fmt.Println(names)

}
