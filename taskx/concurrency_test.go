package taskx

import (
	"context"
	"testing"
	"time"
)

func TestConcurrent(t *testing.T) {
	task := NewConcurrencyTask(10).SetCallback(LogResult)
	for i := 0; i < 20; i++ {
		task.AddTask(&testTask{id: i, ratio: 0.5})
	}
	if err := task.Execute(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestConcurrentAndRetry(t *testing.T) {
	task := NewConcurrencyTask(5).SetCallback(RetryCallback(3, time.Second))
	for i := 0; i < 20; i++ {
		task.AddTask(&testTask{id: i, ratio: 0.5})
	}
	if err := task.Execute(context.Background()); err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Second)
}
