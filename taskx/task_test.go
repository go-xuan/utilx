package taskx

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/go-xuan/utilx/errorx"
)

type testTask struct {
	id    int     // ID
	ratio float64 // 成功率
}

func (t testTask) GetID() string {
	return fmt.Sprintf("test:%d", t.id)
}

func (t testTask) Execute(ctx context.Context) error {
	value := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)
	threshold := int(t.ratio * 100)
	if value <= threshold {
		return nil
	}
	return errorx.Errorf("error: id=%d, value=%d，threshold=%d", t.id, value, threshold)
}

func TestTask(t *testing.T) {
	// 测试任务
	var task Task = &testTask{
		id:    5,
		ratio: 0.5,
	}

	// 执行任务
	if err := task.Execute(t.Context()); err != nil {
		t.Log(err)
	}
}
