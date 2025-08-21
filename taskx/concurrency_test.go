package taskx

import (
	"context"
	"testing"
	"time"
)

func TestConcurrent(t *testing.T) {
	concurrencyTask := NewConcurrencyTask(10).SetCallback(ErrorLogCallback)
	for i := 0; i < 20; i++ {
		concurrencyTask.AddExecute(testTask{id: i, ratio: 0.5}.Execute)
	}
	if err := concurrencyTask.Execute(context.Background()); err != nil {
		t.Log(err)
	}
}

func TestConcurrentAndRetry(t *testing.T) {
	concurrencyTask := NewConcurrencyTask(5).SetCallback(RetryCallback(3, time.Second))
	for i := 0; i < 20; i++ {
		concurrencyTask.AddExecute(testTask{id: i, ratio: 0.5}.Execute)
	}
	if err := concurrencyTask.Execute(context.Background()); err != nil {
		t.Log(err)
	}
	time.Sleep(20 * time.Second)
}
