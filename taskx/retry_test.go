package taskx

import (
	"testing"
	"time"
)

func TestRetryStrategy(t *testing.T) {
	task := testTask{id: 1, ratio: 0.1}
	strategy := NewRetryStrategy(10, time.Second)
	times, err := strategy.Execute(t.Context(), task)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("times", times)
	t.Log("retry strategy execute success")
}

func TestRetry(t *testing.T) {
	task := testTask{id: 1, ratio: 0.1}
	retry := NewRetry(task, 10, time.Second)
	if err := retry.Execute(t.Context()); err != nil {
		t.Log(err)
	}
	t.Log("retry task execute success")
}
