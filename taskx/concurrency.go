package taskx

import (
	"context"
	"sync"

	"github.com/go-xuan/utilx/errorx"
)

// NewConcurrency 创建并发执行器
func NewConcurrency(size int) *Concurrency {
	return &Concurrency{
		size: size,
	}
}

// Concurrency 并发策略
type Concurrency struct {
	size int
}

// Execute 并发执行
func (c *Concurrency) Execute(ctx context.Context, executes []Execute, callback ResultCallback) {
	if len(executes) == 0 {
		return
	}
	if c.size <= 0 {
		c.size = 1
	}
	// 计算结果channel缓冲区大小
	total := len(executes)
	var resultBuffer = total / c.size
	if r := total % c.size; r > 0 {
		resultBuffer = resultBuffer + 1
	}
	var executeCh = make(chan Execute, total)      // 任务管道
	var resultCh = make(chan Result, resultBuffer) // 结果管道
	var wg = &sync.WaitGroup{}

	// 使用单独的协程将所有请求添加进管道，避免阻塞
	go func(list []Execute, ch chan Execute) {
		for _, item := range list {
			ch <- item
		}
		close(ch)
	}(executes, executeCh)

	// 并发处理异步请求
	for i := 0; i < c.size; i++ {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, executeCh chan Execute, resultCh chan Result) {
			defer wg.Done()
			for execute := range executeCh {
				err := execute(ctx)
				resultCh <- &ExecuteResult{execute: execute, error: err}
			}
		}(ctx, wg, executeCh, resultCh)
	}

	// 等待所有请求完成
	go func(wg *sync.WaitGroup, resultCh chan Result) {
		wg.Wait()
		close(resultCh)
	}(wg, resultCh)

	// 从结果管道获取请求结果
	for result := range resultCh {
		callback(result)
	}
	return
}

// NewConcurrencyTask 创建并发执行任务
func NewConcurrencyTask(size int) *ConcurrencyTask {
	return &ConcurrencyTask{
		Concurrency: Concurrency{
			size: size,
		},
		executes: make([]Execute, 0),
	}
}

// ConcurrencyTask 并发执行任务
type ConcurrencyTask struct {
	Concurrency                // 并发策略
	executes    []Execute      // 任务执行函数列表
	callback    ResultCallback // 结果回调函数
}

func (t *ConcurrencyTask) Execute(ctx context.Context) error {
	if len(t.executes) == 0 {
		return errorx.New("concurrency execute is nil")
	}
	if t.callback == nil {
		return errorx.New("result callback is nil")
	}
	// 并发执行子任务
	t.Concurrency.Execute(ctx, t.executes, t.callback)
	return nil
}

// AddExecute 添加执行函数
func (t *ConcurrencyTask) AddExecute(executes ...Execute) *ConcurrencyTask {
	if executes == nil {
		return t
	}
	for _, execute := range executes {
		t.executes = append(t.executes, execute)
	}
	return t
}

// SetCallback 设置结果回调函数
func (t *ConcurrencyTask) SetCallback(callback ResultCallback) *ConcurrencyTask {
	t.callback = callback
	return t
}
