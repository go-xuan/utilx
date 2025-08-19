package taskx

import (
	"context"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

// NewCronTask 创建定时任务
func NewCronTask(name, spec string, task Task) *CronTask {
	return &CronTask{
		name: name,
		spec: spec,
		task: task,
	}
}

// CronTask 定时任务
type CronTask struct {
	entry cron.Entry // 定时任务entry
	name  string     // 任务名
	spec  string     // 定时任务表达式
	task  Task       // 任务
}

// Entry 获取定时任务entry
func (t *CronTask) Entry() cron.Entry {
	return t.entry
}

// Name 获取定时任务名称
func (t *CronTask) Name() string {
	return t.name
}

// Spec 获取定时任务表达式
func (t *CronTask) Spec() string {
	return t.spec
}

// Run 执行定时任务
func (t *CronTask) Run() {
	if err := t.task.Execute(context.Background()); err != nil {
		log.WithField("task_name", t.name).WithError(err).Error("task run failed")
	}
}

// Execute 执行定时任务
func (t *CronTask) Execute(ctx context.Context) error {
	scheduler := NewCronScheduler()
	_ = scheduler.Add(t)
	return scheduler.Execute(ctx)
}

// Wrap 包装任务
func (t *CronTask) Wrap(wrapper Wrapper) {
	t.task = wrapper(t.task)
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
