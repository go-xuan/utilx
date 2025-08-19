package taskx

import (
	"context"
	"fmt"
	"testing"
)

func TestInBatches(t *testing.T) {
	var total = 200
	var s []int
	for i := 0; i < total; i++ {
		s = append(s, i)
	}
	if err := NewSerial(total, 13, func(ctx context.Context, start, end int) error {
		fmt.Printf("%d ==> %d :%v \n", start, end, s[start:end])
		return nil
	}).Execute(context.Background()); err != nil {
		t.Error(err)
	}
}
