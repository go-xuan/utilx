package reflectx

import (
	"fmt"
	"testing"
	"time"
)

func TestGetValue(t *testing.T) {
	v := ValueOf(time.Now())
	fmt.Println(GetValue(v))
	fmt.Println(GetValue(v).String())
	fmt.Println(GetValue(v).Int())
	fmt.Println(GetValue(v).Int64())
	fmt.Println(GetValue(v).Bool())
}
