package taskx

import (
	"context"
	"testing"
)

func TestConcurrent(t *testing.T) {
	concurrency := NewConcurrency(5)
	var tasks []Task
	for i := 0; i < 20; i++ {
		tasks = append(tasks, testTask{id: i, ratio: 0.5})
	}
	concurrency.Execute(t.Context(), tasks, PrintResult)
}

func TestConcurrentScheduler(t *testing.T) {
	concurrency := NewConcurrency(5)
	scheduler := NewConcurrencyScheduler("concurrency", concurrency)
	for i := 0; i < 20; i++ {
		tt := testTask{id: i, ratio: 0.5}
		scheduler.AddTask(tt)
	}
	scheduler.AddResultHook(LogResult)
	if err := scheduler.Execute(context.Background()); err != nil {
		t.Log(err)
	}
}
