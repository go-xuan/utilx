package taskx

import (
	"context"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/idx"
)

// NewCronTask 创建定时任务
func NewCronTask(spec string, execute Execute, name ...string) *CronTask {
	var name_ string
	if len(name) > 0 {
		name_ = name[0]
	} else {
		name_ = idx.SnowFlake().String()
	}
	return &CronTask{name: name_, spec: spec, execute: execute}
}

// CronTask 定时任务
type CronTask struct {
	entry   cron.Entry // 定时任务entry
	name    string     // 任务名
	spec    string     // 定时任务表达式
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
	return t.GetName()
}

func (t *CronTask) Execute(ctx context.Context) error {
	logger := log.WithField("task_type", "cron").WithField("task_name", t.name)
	if err := t.execute(ctx); err != nil {
		logger.WithError(err).Error("cron task execute error")
		return errorx.New("cron task execute error")
	}
	logger.Info("cron task execute success")
	return nil
}

func (t *CronTask) GetName() string {
	return t.name
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
	if t.execute != nil && len(wraps) > 0 {
		for _, wrap := range wraps {
			t.execute = wrap(t.execute)
		}
	}
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
