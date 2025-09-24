package taskx

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/utilx/errorx"
)

// NewQueue 创建队列任务
func NewQueue(id string, execute Execute) *Queue {
	return &Queue{
		id:      id,
		execute: execute,
	}
}

// Queue 队列任务
type Queue struct {
	id      string  // 任务ID
	execute Execute // 任务执行函数
	prev    *Queue  // 指向上一个任务
	next    *Queue  // 指向下一个任务
}

func (q *Queue) GetID() string {
	return q.id
}

func (q *Queue) Execute(ctx context.Context) error {
	logger := log.WithField("task_id", q.GetID()).
		WithField("task_type", "queue")
	if err := q.execute(ctx); err != nil {
		logger.WithField("error", err.Error()).Error("queue task execute error")
		return errorx.New("queue task execute error")
	}
	logger.Info("queue task execute success")
	return nil
}

// HasNext 是否有下一个任务
func (q *Queue) HasNext() bool {
	return q.next != nil
}
