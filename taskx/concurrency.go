package taskx

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

// NewConcurrency 创建并发任务执行器
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
	if size <= 0 {
		size = 8
	}
	return &ConcurrencyStrategy{
		size: size,
	}
}

// ConcurrencyStrategy 并发策略
type ConcurrencyStrategy struct {
	size int // 并发数
}

func (c *ConcurrencyStrategy) Execute(ctx context.Context, tasks []Task, hooks ...ResultHook) {
	total := len(tasks)
	if total == 0 {
		return
	}

	// inCh 缓冲：至少为 c.size，确保发送任务时不会被阻塞
	// 公式：max(c.size, total/c.size)，避免任务数少或并发数大时缓冲过小
	inBuf := total / c.size
	if inBuf < c.size {
		inBuf = c.size
	}
	inCh := make(chan Task, inBuf)
	outCh := make(chan IResult, total) // outCh 全缓冲，确保 worker 不会被结果输出阻塞
	var wg sync.WaitGroup

	// 使用单独的协程将所有任务添加到管道
	go func() {
		defer close(inCh)
		for _, task := range tasks {
			select {
			case <-ctx.Done():
				return
			case inCh <- task:
			}
		}
	}()

	// 并发处理任务
	for i := 0; i < c.size; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range inCh {
				select {
				case <-ctx.Done():
					return
				default:
					err := task.Execute(ctx)
					outCh <- &Result{task: task, error: err}
				}
			}
		}()
	}

	// 等待所有任务完成并关闭结果管道
	go func() {
		wg.Wait()
		close(outCh)
	}()

	// 执行钩子函数
	if len(hooks) > 0 {
		for result := range outCh {
			for _, hook := range hooks {
				hook(result)
			}
		}
	}
}
