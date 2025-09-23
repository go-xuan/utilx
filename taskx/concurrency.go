package taskx

import (
	"context"
	"sync"
)

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
