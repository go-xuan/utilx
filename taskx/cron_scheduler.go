package taskx

import (
	"context"
	"sync"
	"time"

	"github.com/go-xuan/utilx/errorx"
	"github.com/robfig/cron/v3"
)

const (
	initializationStatus = iota // 初始化
	readinessStatus             // 待运行
	runningStatus               // 运行中
	stopStatus                  // 停止
)

// NewCronScheduler 定时任务调度器
func NewCronScheduler(wraps ...Wrap) *CronScheduler {
	scheduler := &CronScheduler{
		mutex:  new(sync.Mutex),
		status: initializationStatus,
		cron: cron.New(
			cron.WithParser(CronParser()),
			cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
			cron.WithLogger(cron.DefaultLogger),
		),
		names:   []string{},
		tasks:   make(map[string]*CronTask),
		wrapper: &Wrapper{wraps: make([]Wrap, 0)},
	}
	scheduler.wrapper.Add(wraps...)
	return scheduler
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

// CronScheduler 定时任务调度器
type CronScheduler struct {
	mutex   *sync.Mutex          // 互斥锁
	cron    *cron.Cron           // corn对象
	status  uint                 // 调度器状态（0-初始化；1-待运行；2-运行中；3-停止）
	names   []string             // 任务名称
	tasks   map[string]*CronTask // 定时任务
	wrapper *Wrapper             // 定时任务包装器
}

// AddWrap 添加包装器
func (s *CronScheduler) AddWrap(wraps ...Wrap) error {
	s.wrapper.Add(wraps...)
	return nil
}

// AddTask 添加定时任务
func (s *CronScheduler) AddTask(task *CronTask) *CronScheduler {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	name, spec := task.name, task.spec

	// 如果已存在同名任务则先移除再新增
	var exist bool
	if old, ok := s.tasks[name]; ok {
		exist = true
		s.cron.Remove(old.entry.ID)
	}

	// 遍历装饰器，对任务执行方法进行包装
	if wraps := s.wrapper.wraps; wraps != nil {
		for _, wrap := range wraps {
			task.Wrap(wrap)
		}
	}

	// 添加定时任务
	entryId, err := s.cron.AddJob(spec, task)
	if err != nil {
		return s
	}
	task.entry = s.cron.Entry(entryId)
	s.tasks[name] = task
	if !exist {
		s.names = append(s.names, name)
	}
	if s.status == initializationStatus {
		s.status = readinessStatus
	}
	return s
}

// Remove 移除定时任务
func (s *CronScheduler) Remove(name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if task, ok := s.tasks[name]; ok {
		s.cron.Remove(task.entry.ID)
		delete(s.tasks, name)
	} else {
		return errorx.New("execute not found: " + name)
	}
	// 当任务清零则状态值归零
	if len(s.tasks) == 0 {
		if s.status == runningStatus {
			s.cron.Stop()
		}
		s.status = initializationStatus
	}
	return nil
}

// Execute 开始调度定时任务
func (s *CronScheduler) Execute(ctx context.Context) error {
	switch s.status {
	case initializationStatus:
		return errorx.New("please add the GetExecute first")
	case runningStatus:
		return errorx.New("the cron scheduler already running")
	default:
		s.cron.Start()
		s.status = runningStatus
		return nil
	}
}

// Stop 停止执行定时任务
func (s *CronScheduler) Stop() error {
	switch s.status {
	case initializationStatus, readinessStatus:
		return errorx.New("the cron scheduler is not running yet")
	case stopStatus:
		return errorx.New("the cron scheduler has stopped")
	default:
		s.cron.Stop()
		s.status = stopStatus
		return nil
	}
}

func (s *CronScheduler) Status() string {
	switch s.status {
	case initializationStatus:
		return "initialization"
	case readinessStatus:
		return "readiness"
	case runningStatus:
		return "running"
	case stopStatus:
		return "stopped"
	default:
		return "unknown"
	}
}

// All 获取所有定时任务
func (s *CronScheduler) All() []*CronTask {
	var tasks []*CronTask
	for _, name := range s.names {
		if task, ok := s.tasks[name]; ok {
			task.entry = s.cron.Entry(task.entry.ID)
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// Get 获取定时任务
func (s *CronScheduler) Get(name string) *CronTask {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if task, ok := s.tasks[name]; ok {
		task.entry = s.cron.Entry(task.entry.ID)
		return task
	}
	return nil
}
