package taskx

import (
	"context"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

// NewCronTask 创建定时任务
func NewCronTask(name, spec string, execute Execute) *CronTask {
	return &CronTask{
		name:    name,
		spec:    spec,
		execute: execute,
	}
}

// CronTask 定时任务
type CronTask struct {
	entry   cron.Entry // 定时任务entry
	spec    string     // 定时任务表达式
	name    string     // 任务名
	execute Execute    // 任务执行函数
}

// Run 执行定时任务
func (t *CronTask) Run() {
	if err := t.execute(context.Background()); err != nil {
		log.WithField("task_name", t.name).WithError(err).Error("execute failed")
	}
}

// GetUnique 获取定时任务名称
func (t *CronTask) GetUnique() string {
	return t.name
}

func (t *CronTask) Execute(ctx context.Context) error {
	return t.execute(ctx)
}

// GetSpec 获取定时任务表达式
func (t *CronTask) GetSpec() string {
	return t.spec
}

// GetEntry 获取定时任务entry
func (t *CronTask) GetEntry() cron.Entry {
	return t.entry
}

// Wrap 包装定时任务执行函数
func (t *CronTask) Wrap(wraps ...Wrap) {
	if t.execute == nil || len(wraps) == 0 {
		return
	}
	execute := t.execute
	for _, wrap := range wraps {
		execute = wrap(execute)
	}
	t.execute = execute
}

// GetMeta 获取定时任务元数据
func (t *CronTask) GetMeta() map[string]string {
	return map[string]string{
		"name": t.name,
		"spec": t.spec,
		"prev": t.entry.Prev.Format("2006-01-02 15:04:05"),
		"next": t.entry.Next.Format("2006-01-02 15:04:05"),
	}
}
