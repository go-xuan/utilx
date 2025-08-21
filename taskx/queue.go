package taskx

import "context"

// QueueTask 队列任务
type QueueTask struct {
	name    string     // 任务名
	execute Execute    // 当前任务执行函数
	prev    *QueueTask // 指向上一个任务
	next    *QueueTask // 指向下一个任务
}

// GetUnique 获取队列任务唯一标识
func (t *QueueTask) GetUnique() string {
	return t.name
}

// Execute 执行任务
func (t *QueueTask) Execute(ctx context.Context) error {
	return t.execute(ctx)
}

// HasNext 是否有下一个任务
func (t *QueueTask) HasNext() bool {
	return t.next != nil
}
