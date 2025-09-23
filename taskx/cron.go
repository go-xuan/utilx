package taskx

import (
	"context"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/utilx/errorx"
)

// NewCron 创建定时任务
func NewCron(id, spec string, execute Execute) *Cron {
	return &Cron{
		id:      id,
		spec:    spec,
		execute: execute,
	}
}

// Cron 定时任务
type Cron struct {
	id      string     // 任务ID
	spec    string     // 定时任务表达式
	execute Execute    // 执行任务函数
	entry   cron.Entry // 定时任务entry
}

// Run 执行定时任务
func (t *Cron) Run() {
	if err := t.Execute(context.Background()); err != nil {
		log.WithField("task_id", t.GetID()).
			WithField("task_type", "cron").
			WithError(err).
			Error("execute failed")
	}
}

// GetID 获取定时任务名称
func (t *Cron) GetID() string {
	return t.id
}

func (t *Cron) Execute(ctx context.Context) error {
	logger := log.WithField("task_id", t.GetID()).
		WithField("task_type", "cron")
	if err := t.execute(ctx); err != nil {
		logger.WithField("error", err.Error()).
			Error("cron task execute error")
		return errorx.New("cron task execute error")
	}
	logger.Info("cron task execute success")
	return nil
}

// GetSpec 获取定时任务表达式
func (t *Cron) GetSpec() string {
	return t.spec
}

// GetEntry 获取定时任务entry
func (t *Cron) GetEntry() cron.Entry {
	return t.entry
}

// Wrap 包装定时任务执行函数
func (t *Cron) Wrap(wraps ...Wrap) {
	execute := t.execute
	if execute != nil && len(wraps) > 0 {
		for _, wrap := range wraps {
			execute = wrap(execute)
		}
	}
	t.execute = execute
}

// GetMeta 获取定时任务元数据
func (t *Cron) GetMeta() map[string]string {
	return map[string]string{
		"id":   t.GetID(),
		"spec": t.GetSpec(),
		"prev": t.GetEntry().Prev.Format("2006-01-02 15:04:05"),
		"next": t.GetEntry().Next.Format("2006-01-02 15:04:05"),
	}
}
