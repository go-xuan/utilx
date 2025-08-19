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
func (c *Concurrency) Execute(ctx context.Context, tasks []Task, callback ResultCallback) {
	if len(tasks) == 0 {
		return
	}
	if c.size <= 0 {
		c.size = 1
	}
	// 计算结果channel缓冲区大小
	totalBuffer := len(tasks)
	var resultBuffer = totalBuffer / c.size
	if r := totalBuffer % c.size; r > 0 {
		resultBuffer = resultBuffer + 1
	}
	var taskCh = make(chan Task, totalBuffer)      // 任务管道
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
		go func(ctx context.Context, wg *sync.WaitGroup, taskCh <-chan Task, resultCh chan<- Result) {
			defer wg.Done()
			for task := range taskCh {
				err := task.Execute(ctx)
				resultCh <- &ConcurrencyResult{task: task, error: err}
			}
		}(ctx, wg, taskCh, resultCh)
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
		tasks: make([]Task, 0),
	}
}

// ConcurrencyTask 并发执行任务
type ConcurrencyTask struct {
	Concurrency                // 并发策略
	tasks       []Task         // 任务实例列表
	callback    ResultCallback // 结果回调函数
}

func (s *ConcurrencyTask) Execute(ctx context.Context) error {
	if len(s.tasks) == 0 {
		return errorx.New("concurrency task is nil")
	}
	if s.callback == nil {
		return errorx.New("result callback is nil")
	}
	// 并发执行子任务
	s.Concurrency.Execute(ctx, s.tasks, s.callback)
	return nil
}

// AddTask 添加子任务
func (s *ConcurrencyTask) AddTask(tasks ...Task) *ConcurrencyTask {
	if tasks == nil {
		return s
	}
	for _, task := range tasks {
		s.tasks = append(s.tasks, task)
	}
	return s

}

// SetCallback 设置结果回调函数
func (s *ConcurrencyTask) SetCallback(callback ResultCallback) *ConcurrencyTask {
	s.callback = callback
	return s
}

// ConcurrencyResult 并发执行结果
type ConcurrencyResult struct {
	task  Task
	error error
}

func (r *ConcurrencyResult) Task() Task {
	return r.task
}

func (r *ConcurrencyResult) Error() error {
	return r.error
}
