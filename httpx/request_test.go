package httpx

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	resp, err := NewRequest(http.MethodGet, "http://localhost:8080/ping").Debug().Send()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.StatusOK())
}
