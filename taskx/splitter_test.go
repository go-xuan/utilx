package taskx

import (
	"context"
	"fmt"
	"testing"
)

func TestSplitter(t *testing.T) {
	var total = 200
	var s []int
	for i := 1; i <= total; i++ {
		s = append(s, i)
	}
	if err := NewSplitter(17).Execute(t.Context(), total, func(ctx context.Context, start, end, batch int) error {
		fmt.Printf("%d ==> [%d:%d] ==> %v \n", batch, start, end, s[start:end])
		return nil
	}); err != nil {
		t.Log(err)
	}
}

func TestSplitterTask(t *testing.T) {
	var total = 200
	var s []int
	for i := 1; i <= total; i++ {
		s = append(s, i)
	}
	task := NewSplitterTask[int](17)
	if err := task.SetList(s).SetExecute(func(ctx context.Context, list []int) error {
		fmt.Printf("%v \n", list)
		return nil
	}).Execute(t.Context()); err != nil {
		t.Log(err)
	}
}
