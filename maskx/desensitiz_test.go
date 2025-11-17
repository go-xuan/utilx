package maskx

import (
	"fmt"
	"testing"

	"github.com/go-xuan/utilx/randx"
)

func TestDesensitize(t *testing.T) {
	phone := randx.Phone()
	name := randx.Name()
	email := randx.Email()
	fmt.Println(phone, "==>", Phone.Desensitize(phone))
	fmt.Println(name, "==>", Name.Desensitize(name))
	fmt.Println(email, "==>", Email.Desensitize(email))
}
