package taskx

import (
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	if _, err := NewRetry(10, time.Second).Execute(t.Context(), testTask{id: 1, ratio: 0.1}); err != nil {
		t.Log(err)
	}
}

func TestRetryTask(t *testing.T) {
	// 执行
	if err := NewRetryTask(testTask{id: 1, ratio: 0.1}, 10, time.Second).Execute(t.Context()); err != nil {
		t.Log(err)
		return
	}
	t.Log("retry task execute success")
}
