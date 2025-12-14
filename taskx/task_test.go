package taskx

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/funcx"
)

func function(id int, ratio float64) funcx.X {
	return func(ctx context.Context) error {
		value := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)
		threshold := int(ratio * 100)
		if value <= threshold {
			return nil
		}
		return errorx.Sprintf("error: id=%d, value=%d，threshold=%d", id, value, threshold)
	}
}

type testTask struct {
	id    int     // ID
	ratio float64 // 成功率
}

func (t testTask) GetID() string {
	return fmt.Sprintf("test—%d", t.id)
}

func (t testTask) Execute(ctx context.Context) error {
	return function(t.id, t.ratio)(ctx)
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
