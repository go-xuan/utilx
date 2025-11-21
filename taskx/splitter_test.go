package taskx

import (
	"context"
	"fmt"
	"testing"
)

func TestSplitterStrategy(t *testing.T) {
	var total = 200
	var tasks []int
	for i := 1; i <= total; i++ {
		tasks = append(tasks, i)
	}
	execute := func(ctx context.Context, start, end, batch int) error {
		fmt.Printf("%d ==> [%d:%d] ==> %v \n", batch, start, end, tasks[start:end])
		return nil
	}
	if err := NewSplitterStrategy(17).Execute(t.Context(), total, execute); err != nil {
		t.Log(err)
	}
}

func TestSplitter(t *testing.T) {
	var tasks []testTask
	for i := 1; i <= 200; i++ {
		tasks = append(tasks, testTask{id: i, ratio: 0.5})
	}
	execute := func(ctx context.Context, tasks []testTask) error {
		// 模拟执行任务
		ids := make([]int, 0, len(tasks))
		for _, task := range tasks {
			ids = append(ids, task.id)
		}
		fmt.Printf("%v \n", ids)
		return nil
	}

	scheduler := NewSplitter[testTask](17).SetBatchExecute(execute).AddTask(tasks...)
	if err := scheduler.Execute(t.Context()); err != nil {
		t.Log(err)
	}
}
