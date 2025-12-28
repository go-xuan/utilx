package taskx

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/funcx"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

const (
	initializationStatus = iota // 初始化
	readinessStatus             // 就绪
	runningStatus               // 运行中
	stoppedStatus               // 停止
)

// NewCronScheduler 定时任务调度器
func NewCronScheduler(id string, opts ...cron.Option) *CronScheduler {
	return &CronScheduler{
		id:     id,
		mutex:  new(sync.Mutex),
		cron:   cron.New(opts...),
		status: initializationStatus,
		names:  []string{},
		jobs:   make(map[string]*CronJob),
		wraps:  make([]funcx.Wrap[funcx.X], 0),
	}
}

// CronScheduler 定时任务调度器
type CronScheduler struct {
	id     string                // 定时任务调度器ID
	mutex  *sync.Mutex           // 互斥锁
	cron   *cron.Cron            // corn对象
	status uint                  // 调度器状态（0-初始化；1-待运行；2-运行中；3-停止）
	names  []string              // 任务entryID列表
	jobs   map[string]*CronJob   // 定时任务
	wraps  []funcx.Wrap[funcx.X] // 定时任务执行函数包装器
}

func (c *CronScheduler) GetID() string {
	return c.id
}

func (c *CronScheduler) Execute(_ context.Context) error {
	logger := log.WithField("task_id", c.GetID()).WithField("task_type", "cron_scheduler")
	if err := c.Start(); err != nil {
		logger.WithError(err).Error("cron scheduler start failed")
		return errorx.Wrap(err, "cron scheduler start failed")
	}
	logger.Info("cron scheduler start")
	return nil
}

// Start 启动定时任务调度器
func (c *CronScheduler) Start() error {
	switch c.status {
	case initializationStatus:
		return errorx.New("cron scheduler initialization")
	case runningStatus:
		return errorx.New("cron scheduler already running")
	default:
		c.cron.Start()
		c.status = runningStatus
		return nil
	}
}

// AddWrap 添加包装器
func (c *CronScheduler) AddWrap(wraps ...funcx.Wrap[funcx.X]) *CronScheduler {
	c.wraps = append(c.wraps, wraps...)
	return c
}

// AddJob 添加定时任务
func (c *CronScheduler) AddJob(name, spec string, function funcx.X) *CronScheduler {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 创建定时任务
	job := &CronJob{name: name, spec: spec}

	// 如果已存在同名任务则先移除再新增
	var exist bool
	if old, ok := c.jobs[name]; ok {
		exist = true
		c.cron.Remove(old.GetEntry().ID)
		job.createTime = old.createTime
		job.times = old.times
	}

	// 遍历装饰器，对任务执行函数进行包装
	job.function = funcx.WrapX(function, c.wraps...)

	// 添加定时任务
	entryId, err := c.cron.AddJob(spec, job)
	if err != nil {
		return c
	}

	if !exist {
		job.createTime = time.Now()
		c.names = append(c.names, name)
	}
	job.updateTime = time.Now()
	job.entry = c.cron.Entry(entryId)

	c.jobs[name] = job
	if c.status == initializationStatus {
		c.status = readinessStatus
	}
	return c
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
func (c *CronScheduler) All() []*CronJob {
	var jobs []*CronJob
	for _, name := range c.names {
		if job, ok := c.jobs[name]; ok {
			job.entry = c.cron.Entry(job.GetEntryID())
			jobs = append(jobs, job)
		}
	}
	return jobs
}

// Get 获取定时任务
func (c *CronScheduler) Get(name string) *CronJob {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if job, ok := c.jobs[name]; ok {
		job.entry = c.cron.Entry(job.GetEntryID())
		return job
	}
	return nil
}

// Remove 移除定时任务
func (c *CronScheduler) Remove(name string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if job, ok := c.jobs[name]; ok {
		c.cron.Remove(job.GetEntryID())
		delete(c.jobs, name)
	} else {
		return errorx.Sprintf("task not found: %s", name)
	}
	// 当任务清零则回归初始化状态
	if len(c.jobs) == 0 {
		if c.status == runningStatus {
			c.cron.Stop()
		}
		c.status = initializationStatus
	}
	return nil
}

// ParseCronSpec 解cron析表达式，计算当前时间和下次执行时间的时间差
func ParseCronSpec(parser cron.Parser, spec string) time.Duration {
	if schedule, err := parser.Parse(spec); err == nil {
		now := time.Now()
		next := schedule.Next(now)
		return next.Sub(now)
	}
	return -1
}

// CronJob 定时任务
type CronJob struct {
	name       string     // 定时任务名称，唯一标识
	spec       string     // 定时任务表达式
	times      int        // 执行次数
	function   funcx.X    // 执行函数
	entry      cron.Entry // 定时任务entry
	createTime time.Time  // 创建时间
	updateTime time.Time  // 更新时间
}

// Run 执行定时任务
func (t *CronJob) Run() {
	_ = t.Execute(context.Background())
	t.times++
}

func (t *CronJob) GetID() string {
	return t.name
}

func (t *CronJob) Execute(ctx context.Context) error {
	logger := log.WithField("task_id", t.GetID())
	if err := t.function(ctx); err != nil {
		logger.WithError(err).Error("cron job execute error")
		return errorx.Wrap(err, "cron job execute error")
	}
	logger.Info("cron job execute success")
	return nil
}

// GetEntryID 获取定时任务entry ID
func (t *CronJob) GetEntryID() cron.EntryID {
	return t.entry.ID
}

// GetSpec 获取定时任务表达式
func (t *CronJob) GetSpec() string {
	return t.spec
}

// GetTimes 获取定时任务执行次数
func (t *CronJob) GetTimes() int {
	return t.times
}

// GetEntry 获取定时任务entry
func (t *CronJob) GetEntry() cron.Entry {
	return t.entry
}

// GetCreateTime 获取定时任务创建时间
func (t *CronJob) GetCreateTime() time.Time {
	return t.createTime
}

// GetUpdateTime 获取定时任务更新时间
func (t *CronJob) GetUpdateTime() time.Time {
	return t.updateTime
}

// GetMeta 获取定时任务元数据
func (t *CronJob) GetMeta() map[string]string {
	layout := "2006-01-02 15:04:05"
	return map[string]string{
		"id":          t.GetID(),
		"spec":        t.GetSpec(),
		"prev":        t.GetEntry().Prev.Format(layout),
		"next":        t.GetEntry().Next.Format(layout),
		"times":       fmt.Sprintf("%d", t.times),
		"create_time": t.createTime.Format(layout),
		"update_time": t.updateTime.Format(layout),
	}
}

// NewCronLogger 创建cron日志记录器
func NewCronLogger() *CronLogger {
	return &CronLogger{
		logger: log.WithField("logger", "cron"),
	}
}

type CronLogger struct {
	logger *log.Entry
}

func (l *CronLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.WithFields(logFields(keysAndValues...)).Info(msg)
}

func (l *CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.WithError(err).WithFields(logFields(keysAndValues...)).Error(msg)
}

func logFields(keysValues ...interface{}) map[string]interface{} {
	var fields = make(map[string]interface{})
	length := len(keysValues)
	if length == 0 {
		return fields
	}
	for i := 0; i < length/2; i++ {
		if ki, vi := i*2, i*2+1; vi < length {
			key := fmt.Sprintf("%v", keysValues[ki])
			value := fmt.Sprintf("%v", keysValues[vi])
			fields[key] = value
		}
	}
	return fields
}
