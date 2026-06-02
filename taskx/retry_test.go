package taskx

import (
	"testing"
	"time"
)

func TestRetryStrategy(t *testing.T) {
	task := testTask{id: 1, ratio: 0.1}
	err := NewRetryStrategy(6, 5*time.Second, 1.5).Execute(t.Context(), task)
	if err != nil {
		t.Log(err)
		return
	}
}

func TestRetry(t *testing.T) {
	task := testTask{id: 1, ratio: 0.1}
	retry := NewRetry(task, DefaultRetryStrategy())
	if err := retry.Execute(t.Context()); err != nil {
		t.Log(err)
	}
	t.Log("retry task execute success")
}
