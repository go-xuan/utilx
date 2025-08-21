package httpx

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	if resp, err := NewRequest(http.MethodGet, "http://localhost:3456/tools/demo").Debug().Send(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusOK())
	}
}
