package taskx

import (
	"context"
	"fmt"
	"testing"
)

func TestSplitter(t *testing.T) {
	var total = 200
	var tasks []int
	for i := 1; i <= total; i++ {
		tasks = append(tasks, i)
	}
	splitter := NewSplitter(17)
	
	execute := func(ctx context.Context, start, end, batch int) error {
		fmt.Printf("%d ==> [%d:%d] ==> %v \n", batch, start, end, tasks[start:end])
		return nil
	}

	if err := splitter.Execute(t.Context(), total, execute); err != nil {
		t.Log(err)
	}
}

func TestSplitterScheduler(t *testing.T) {
	var total = 200
	var tasks []int
	for i := 1; i <= total; i++ {
		tasks = append(tasks, i)
	}
	execute := func(ctx context.Context, ints []int) error {
		// 打印数字
		fmt.Printf("%v \n", ints)
		return nil
	}

	scheduler := NewSplitterScheduler("splitter", NewSplitter(17), tasks...).SetBatchExecute(execute)

	if err := scheduler.Execute(t.Context()); err != nil {
		t.Log(err)
	}
}
