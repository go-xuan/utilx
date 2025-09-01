package taskx

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/idx"
)

// NewQueueTask 创建队列任务
func NewQueueTask(execute Execute, name ...string) *QueueTask {
	var name_ string
	if len(name) > 0 {
		name_ = name[0]
	} else {
		name_ = idx.SnowFlake().String()
	}
	return &QueueTask{name: name_, execute: execute}
}

// QueueTask 队列任务
type QueueTask struct {
	name    string     // 任务名
	execute Execute    // 任务执行函数
	prev    *QueueTask // 指向上一个任务
	next    *QueueTask // 指向下一个任务
}

func (t *QueueTask) GetUnique() string {
	return t.GetName()
}

func (t *QueueTask) Execute(ctx context.Context) error {
	logger := log.WithField("task_type", "queue").WithField("task_name", t.name)
	if err := t.execute(ctx); err != nil {
		logger.WithError(err).Error("queue task execute error")
		return errorx.New("queue task execute error")
	}
	logger.Info("queue task execute success")
	return t.execute(ctx)
}

func (t *QueueTask) GetName() string {
	return t.name
}

// HasNext 是否有下一个任务
func (t *QueueTask) HasNext() bool {
	return t.next != nil
}
