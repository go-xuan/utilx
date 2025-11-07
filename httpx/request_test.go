package httpx

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	if resp, err := NewRequest(http.MethodGet, "http://localhost:8080/ping").Debug().Send(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusOK())
	}
}
