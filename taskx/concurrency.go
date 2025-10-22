package taskx

import (
	"context"
	"sync"

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

// NewConcurrency 创建并发策略
func NewConcurrency(size int) *Concurrency {
	return &Concurrency{
		size: size,
	}
}

// Concurrency 并发策略
type Concurrency struct {
	size int
}

func (c *Concurrency) Execute(ctx context.Context, tasks []Task, hooks ...ResultHook) {
	if len(tasks) == 0 {
		return
	}
	if c.size <= 0 {
		c.size = 1
	}
	// 计算结果channel缓冲区大小
	total := len(tasks)
	var resultBuffer = total / c.size
	if r := total % c.size; r > 0 {
		resultBuffer = resultBuffer + 1
	}
	var taskCh = make(chan Task, total)            // 任务管道
	var resultCh = make(chan Result, resultBuffer) // 结果管道
	var wg = &sync.WaitGroup{}

	// 使用单独的协程将所有请求添加进管道，避免阻塞
	go func(tasks []Task, taskCh chan Task) {
		for _, task := range tasks {
			taskCh <- task
		}
		close(taskCh)
	}(tasks, taskCh)

	// 并发处理异步请求
	for i := 0; i < c.size; i++ {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, taskCh chan Task, resultCh chan Result) {
			defer wg.Done()
			for task := range taskCh {
				err := task.Execute(ctx)
				resultCh <- &TaskResult{
					task:  task,
					error: err,
				}
			}
		}(ctx, wg, taskCh, resultCh)
	}

	// 等待所有请求完成
	go func(wg *sync.WaitGroup, resultCh chan Result) {
		wg.Wait()
		close(resultCh)
	}(wg, resultCh)

	// 没有钩子函数直接返回
	if len(hooks) == 0 {
		return
	}

	// 执行钩子函数
	for result := range resultCh {
		for _, hook := range hooks {
			hook(result)
		}
	}
	return
}
