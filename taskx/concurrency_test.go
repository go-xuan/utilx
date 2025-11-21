package taskx

import (
	"context"
	"testing"
)

func TestConcurrency(t *testing.T) {
	concurrency := NewConcurrency(5)
	for i := 0; i < 20; i++ {
		tt := testTask{id: i, ratio: 0.5}
		concurrency.AddTask(tt)
	}
	concurrency.AddResultHook(LogResult)
	if err := concurrency.Execute(context.Background()); err != nil {
		t.Log(err)
	}
}
