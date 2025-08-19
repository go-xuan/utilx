package taskx

import (
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	if err := NewRetryTask(10, time.Second).
		AddTask(&testTask{id: 1, ratio: 0.2}).
		Execute(t.Context()); err != nil {
		t.Error(err)
	} else {
		t.Log("task task success")
	}
	time.Sleep(11 * time.Second)
	return
}
