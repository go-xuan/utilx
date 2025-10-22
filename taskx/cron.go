package taskx

import (
	"context"
	"sync"
	"time"

	"github.com/go-xuan/utilx/errorx"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

const (
	initializationStatus = iota // 初始化
	readinessStatus             // 待运行
	runningStatus               // 运行中
	stoppedStatus               // 停止
)

// NewCronScheduler 定时任务调度器
func NewCronScheduler(id string) *CronScheduler {
	return &CronScheduler{
		id:     id,
		mutex:  new(sync.Mutex),
		status: initializationStatus,
		cron: cron.New(
			cron.WithParser(CronParser()),
			cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
			cron.WithLogger(cron.DefaultLogger),
		),
		ids:   []string{},
		tasks: make(map[string]*CronTask),
		wraps: make([]Wrap, 0),
	}
}

// CronScheduler 定时任务调度器
type CronScheduler struct {
	id     string               // 定时任务调度器ID
	mutex  *sync.Mutex          // 互斥锁
	cron   *cron.Cron           // corn对象
	status uint                 // 调度器状态（0-初始化；1-待运行；2-运行中；3-停止）
	ids    []string             // 任务名称
	tasks  map[string]*CronTask // 定时任务
	wraps  []Wrap               // 定时任务包装器
}

func (c *CronScheduler) GetID() string {
	return c.id
}

func (c *CronScheduler) Execute(ctx context.Context) error {
	logger := log.WithField("task_id", c.GetID()).
		WithField("task_type", "cron_scheduler")
	switch c.status {
	case initializationStatus:
		logger.Warn("cron scheduler initialization")
		return errorx.New("cron scheduler initialization")
	case runningStatus:
		logger.Warn("cron scheduler already running")
		return errorx.New("cron scheduler already running")
	default:
		c.cron.Start()
		c.status = runningStatus
		logger.Info("cron scheduler running")
		return nil
	}
}

// AddWrap 添加包装器
func (c *CronScheduler) AddWrap(wraps ...Wrap) *CronScheduler {
	c.wraps = append(c.wraps, wraps...)
	return c
}

// AddCronTask 添加定时任务
func (c *CronScheduler) AddCronTask(task *CronTask) *CronScheduler {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	id, spec := task.GetID(), task.GetSpec()

	// 如果已存在同名任务则先移除再新增
	var exist bool
	if old, ok := c.tasks[id]; ok {
		exist = true
		c.cron.Remove(old.entry.ID)
	}

	// 遍历装饰器，对任务执行函数进行包装
	task.Wrap(c.wraps...)

	// 添加定时任务
	entryId, err := c.cron.AddJob(spec, task)
	if err != nil {
		return c
	}
	task.entry = c.cron.Entry(entryId)
	c.tasks[id] = task
	if !exist {
		c.ids = append(c.ids, id)
	}
	if c.status == initializationStatus {
		c.status = readinessStatus
	}
	return c
}

// Remove 移除定时任务
func (c *CronScheduler) Remove(id string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if task, ok := c.tasks[id]; ok {
		c.cron.Remove(task.entry.ID)
		delete(c.tasks, id)
	} else {
		return errorx.New("execute not found: " + id)
	}
	// 当任务清零则状态值归零
	if len(c.tasks) == 0 {
		if c.status == runningStatus {
			c.cron.Stop()
		}
		c.status = initializationStatus
	}
	return nil
}

// Stop 停止执行定时任务
func (c *CronScheduler) Stop() error {
	switch c.status {
	case initializationStatus, readinessStatus:
		return errorx.New("cron scheduler is not running yet")
	case stoppedStatus:
		return errorx.New("cron scheduler has stopped")
	default:
		c.cron.Stop()
		c.status = stoppedStatus
		return nil
	}
}

func (c *CronScheduler) Status() string {
	switch c.status {
	case initializationStatus:
		return "initialization"
	case readinessStatus:
		return "readiness"
	case runningStatus:
		return "running"
	case stoppedStatus:
		return "stopped"
	default:
		return "unknown"
	}
}

// All 获取所有定时任务
func (c *CronScheduler) All() []*CronTask {
	var tasks []*CronTask
	for _, id := range c.ids {
		if task, ok := c.tasks[id]; ok {
			task.entry = c.cron.Entry(task.entry.ID)
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// Get 获取定时任务
func (c *CronScheduler) Get(id string) *CronTask {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if task, ok := c.tasks[id]; ok {
		task.entry = c.cron.Entry(task.entry.ID)
		return task
	}
	return nil
}

// CronParser 默认的定时任务表达式解析器
func CronParser() cron.Parser {
	return cron.NewParser(
		cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
}

// ParseDurationBySpec 解析表达式，计算当前时间和下次执行时间的时间差
func ParseDurationBySpec(spec string) time.Duration {
	if schedule, err := CronParser().Parse(spec); err == nil {
		var now = time.Now()
		return schedule.Next(now).Sub(now)
	}
	return time.Duration(-1)
}

// NewCronTask 创建定时任务
func NewCronTask(id, spec string, execute Execute) *CronTask {
	return &CronTask{
		id:      id,
		spec:    spec,
		execute: execute,
	}
}

// CronTask 定时任务
type CronTask struct {
	id      string     // 任务ID
	spec    string     // 定时任务表达式
	execute Execute    // 执行任务函数
	entry   cron.Entry // 定时任务entry
}

// Run 执行定时任务
func (t *CronTask) Run() {
	if err := t.Execute(context.Background()); err != nil {
		log.WithField("id", t.GetID()).
			WithError(err).Error("cron task execute error")
	}
}

func (t *CronTask) GetID() string {
	return t.id
}

func (t *CronTask) Execute(ctx context.Context) error {
	logger := log.WithField("id", t.GetID())
	if err := t.execute(ctx); err != nil {
		logger.WithField("error", err.Error()).
			Error("cron task execute error")
		return errorx.New("cron task execute error")
	}
	logger.Info("cron task execute success")
	return nil
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
	execute := t.execute
	if execute != nil && len(wraps) > 0 {
		for _, wrap := range wraps {
			execute = wrap(execute)
		}
	}
	t.execute = execute
}

// GetMeta 获取定时任务元数据
func (t *CronTask) GetMeta() map[string]string {
	return map[string]string{
		"id":   t.GetID(),
		"spec": t.GetSpec(),
		"prev": t.GetEntry().Prev.Format("2006-01-02 15:04:05"),
		"next": t.GetEntry().Next.Format("2006-01-02 15:04:05"),
	}
}
