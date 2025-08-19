package taskx

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/go-xuan/utilx/errorx"
)

type testTask struct {
	id    int     // ID
	ratio float64 // 成功率
}

func (t *testTask) Execute(ctx context.Context) error {
	value := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)
	threshold := int(t.ratio * 100)
	if value <= threshold {
		return nil
	}
	return errorx.Errorf("error: id=%d, value=%d，threshold=%d", t.id, value, threshold)
}

func TestTask(t *testing.T) {
	task := &testTask{5, 0.5}
	if err := task.Execute(context.Background()); err != nil {
		t.Log(err)
	}
	return
}
