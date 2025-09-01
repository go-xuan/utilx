package taskx

import (
	"context"
	"sync"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/marshalx"
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
				resultCh <- &TaskResult{unique: task.GetUnique(), execute: task.Execute, error: err}
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

// NewConcurrencyTask 创建并发执行任务
func NewConcurrencyTask(size int) *ConcurrencyTask {
	return &ConcurrencyTask{
		concurrency: NewConcurrency(size),
		tasks:       make([]Task, 0),
		hooks:       make([]ResultHook, 0),
	}
}

// ConcurrencyTask 并发执行任务
type ConcurrencyTask struct {
	concurrency *Concurrency // 并发策略
	tasks       []Task       // 任务执行函数列表
	hooks       []ResultHook // 结果回调函数
}

func (t *ConcurrencyTask) Execute(ctx context.Context) error {
	if len(t.tasks) == 0 {
		return errorx.New("concurrency tasks is nil")
	}
	// 并发执行子任务
	t.concurrency.Execute(ctx, t.tasks, t.hooks...)
	return nil
}

// AddTask 添加执行函数
func (t *ConcurrencyTask) AddTask(tasks ...Task) *ConcurrencyTask {
	if len(tasks) > 0 {
		t.tasks = append(t.tasks, tasks...)
	}
	return t
}

// AddResultHook 添加结果回调函数
func (t *ConcurrencyTask) AddResultHook(hooks ...ResultHook) *ConcurrencyTask {
	if len(hooks) > 0 {
		t.hooks = append(t.hooks, hooks...)
	}
	return t
}

// ConcurrencyErrorCollect 并发任务错误收集
type ConcurrencyErrorCollect struct {
	Count   int      `json:"count"`   // 失败任务数量
	Uniques []string `json:"uniques"` // 失败任务列表
}

// Save 并发任务错误信息保存
func (c *ConcurrencyErrorCollect) Save(path string) error {
	if c != nil && path != "" && c.Count > 0 {
		return marshalx.Apply(path).Write(path, c)
	}
	return nil
}

// ResultHook 并发任务错误收集钩子
func (c *ConcurrencyErrorCollect) ResultHook(result Result) {
	if err := result.GetError(); err != nil {
		c.Count++
		c.Uniques = append(c.Uniques, result.GetUnique())
	}
}
