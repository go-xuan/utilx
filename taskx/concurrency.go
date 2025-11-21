package taskx

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

// NewConcurrency 创建
func NewConcurrency(size int) *Concurrency {
	return &Concurrency{
		tasks:    make([]Task, 0),
		hooks:    make([]ResultHook, 0),
		strategy: NewConcurrencyStrategy(size),
	}
}

// Concurrency 并发执行任务
type Concurrency struct {
	tasks    []Task               // 任务执行函数列表
	hooks    []ResultHook         // 结果回调函数
	strategy *ConcurrencyStrategy // 并发策略
}

func (t *Concurrency) GetID() string {
	return fmt.Sprintf("concurrency@%d", t.strategy.size)
}

func (t *Concurrency) Execute(ctx context.Context) error {
	logger := log.WithField("task_id", t.GetID())
	if len(t.tasks) == 0 {
		logger.Error("concurrency scheduler tasks is nil")
		return errorx.New("concurrency scheduler tasks is nil")
	} else if t.strategy == nil {
		logger.Error("concurrency strategy is nil")
		return errorx.New("concurrency strategy is nil")
	}
	// 并发执行子任务
	t.strategy.Execute(ctx, t.tasks, t.hooks...)
	logger.Infof("concurrency execute finished")
	return nil
}

// AddTask 添加执行函数
func (t *Concurrency) AddTask(tasks ...Task) *Concurrency {
	if len(tasks) > 0 {
		t.tasks = append(t.tasks, tasks...)
	}
	return t
}

// AddResultHook 添加结果回调函数
func (t *Concurrency) AddResultHook(hooks ...ResultHook) *Concurrency {
	if len(hooks) > 0 {
		t.hooks = append(t.hooks, hooks...)
	}
	return t
}

// NewConcurrencyStrategy 创建并发策略
func NewConcurrencyStrategy(size int) *ConcurrencyStrategy {
	return &ConcurrencyStrategy{
		size: size,
	}
}

// ConcurrencyStrategy 并发策略
type ConcurrencyStrategy struct {
	size int // 并发数
}

func (c *ConcurrencyStrategy) Execute(ctx context.Context, tasks []Task, hooks ...ResultHook) {
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
	var taskCh = make(chan Task, total)             // 任务管道
	var resultCh = make(chan IResult, resultBuffer) // 结果管道
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
		go func(ctx context.Context, wg *sync.WaitGroup, taskCh chan Task, resultCh chan IResult) {
			defer wg.Done()
			for task := range taskCh {
				err := task.Execute(ctx)
				resultCh <- &Result{
					task:  task,
					error: err,
				}
			}
		}(ctx, wg, taskCh, resultCh)
	}

	// 等待所有请求完成
	go func(wg *sync.WaitGroup, resultCh chan IResult) {
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
