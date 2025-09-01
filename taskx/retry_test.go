package taskx

import (
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	if err := NewRetryTask("test_retry", NewRetry(10, time.Second)).
		AddExecute(testTask{id: 1, ratio: 0.1}.Execute).
		Execute(t.Context()); err != nil {
		t.Log(err)
	} else {
		t.Log("execute success")
	}
}
