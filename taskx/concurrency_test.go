package taskx

import (
	"context"
	"testing"
	"time"
)

func TestConcurrent(t *testing.T) {
	concurrencyTask := NewConcurrencyTask(10).AddResultHook(ErrorLogHook)
	for i := 0; i < 20; i++ {
		concurrencyTask.AddTask(testTask{id: i, ratio: 0.5})
	}
	if err := concurrencyTask.Execute(context.Background()); err != nil {
		t.Log(err)
	}
}

func TestConcurrentAndRetry(t *testing.T) {
	concurrencyTask := NewConcurrencyTask(5).AddResultHook(NewRetry(3, time.Second).ResultHook)
	for i := 0; i < 20; i++ {
		concurrencyTask.AddTask(testTask{id: i, ratio: 0.5})
	}
	if err := concurrencyTask.Execute(context.Background()); err != nil {
		t.Log(err)
	}
	time.Sleep(20 * time.Second)
}
