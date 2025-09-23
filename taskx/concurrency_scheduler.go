package taskx

import (
	"context"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

// NewConcurrencyScheduler 创建并发执行任务
func NewConcurrencyScheduler(id string, concurrency *Concurrency) *ConcurrencyScheduler {
	return &ConcurrencyScheduler{
		id:          id,
		concurrency: concurrency,
		tasks:       make([]Task, 0),
		hooks:       make([]ResultHook, 0),
	}
}

// ConcurrencyScheduler 并发执行任务
type ConcurrencyScheduler struct {
	id          string       // 任务ID
	concurrency *Concurrency // 并发策略
	tasks       []Task       // 任务执行函数列表
	hooks       []ResultHook // 结果回调函数
}

func (t *ConcurrencyScheduler) GetID() string {
	return t.id
}

func (t *ConcurrencyScheduler) Execute(ctx context.Context) error {
	logger := log.WithField("task_id", t.GetID()).
		WithField("task_type", "concurrency_scheduler")
	if len(t.tasks) == 0 {
		logger.Error("concurrency scheduler tasks is nil")
		return errorx.New("concurrency scheduler tasks is nil")
	}
	if t.concurrency == nil {
		return errorx.New("concurrency is nil")
	}
	// 并发执行子任务
	t.concurrency.Execute(ctx, t.tasks, t.hooks...)
	logger.Infof("concurrency scheduler execute finished")
	return nil
}

// AddTask 添加执行函数
func (t *ConcurrencyScheduler) AddTask(tasks ...Task) *ConcurrencyScheduler {
	if len(tasks) > 0 {
		t.tasks = append(t.tasks, tasks...)
	}
	return t
}

// AddResultHook 添加结果回调函数
func (t *ConcurrencyScheduler) AddResultHook(hooks ...ResultHook) *ConcurrencyScheduler {
	if len(hooks) > 0 {
		t.hooks = append(t.hooks, hooks...)
	}
	return t
}
