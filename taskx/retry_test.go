package taskx

import (
	"testing"
	"time"
)

func TestRetryScheduler(t *testing.T) {
	retry := NewRetry(10, time.Second) // 重试策略
	tt := testTask{id: 1, ratio: 0.1}  // 任务

	// 执行
	if err := NewRetryScheduler(tt, retry).Execute(t.Context()); err != nil {
		t.Log(err)
	} else {
		t.Log("retry execute success")
	}
}
