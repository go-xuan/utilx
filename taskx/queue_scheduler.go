package taskx

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/utilx/errorx"
)

// NewQueueScheduler 队列任务处理调度器
func NewQueueScheduler(id string) *QueueScheduler {
	return &QueueScheduler{
		id:    id,
		mutex: new(sync.Mutex),
		tasks: make(map[string]*Queue),
	}
}

// QueueScheduler 队列任务处理调度器
type QueueScheduler struct {
	id    string            // 队列任务处理调度器ID
	mutex *sync.Mutex       // 锁
	head  *Queue            // 头部任务
	tail  *Queue            // 尾部任务
	tasks map[string]*Queue // 任务列表
}

func (q *QueueScheduler) GetID() string {
	return q.id
}

func (q *QueueScheduler) Execute(ctx context.Context) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	logger := log.WithField("task_id", q.GetID()).
		WithField("task_type", "queue_scheduler")

	// 检查队列是否为空
	if q.head == nil || len(q.tasks) == 0 {
		logger.Info("queue tasks is empty")
		return nil
	}

	current := q.head
	for current != nil {
		logger = logger.WithField("curr_task_id", current.GetID())
		if current.next != nil {
			logger = logger.WithField("next_task_id", current.next.GetID())
		}

		// 执行当前任务
		if err := current.Execute(ctx); err != nil {
			logger.WithField("error", err.Error()).Error("queue task schedule error")
			return errorx.Wrap(err, "queue task schedule error")
		}
		logger.Info("queue task schedule success")

		// 从任务列表中删除当前任务并更新当前任务指针
		delete(q.tasks, current.id)
		current = current.next
	}
	return nil
}

// Add 新增任务（默认尾插）
func (q *QueueScheduler) Add(id string, execute Execute) {
	q.AddTail(id, execute)
}

// AddTail 尾插（当前新增任务添加到队列末尾）
func (q *QueueScheduler) AddTail(id string, execute Execute) {
	logger := log.WithField("add_task_id", id).
		WithField("add_position", "queue_tail")
	if id == "" || execute == nil {
		logger.Error("tasks id is empty or tasks execute is nil")
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.cover(id, execute) {
		logger.Error("queue add tail failed: tasks already exists")
		return
	}

	newTask := NewQueue(id, execute)
	if q.tail != nil {
		// 队列不为空，将新任务添加到尾部
		newTask.prev = q.tail
		q.tail.next = newTask
		q.tail = newTask
	} else {
		// 队列为空，新任务既是头也是尾
		q.head = newTask
		q.tail = newTask
	}

	q.tasks[id] = newTask
	logger.Info("queue add tail success")
}

// AddHead 头插（当前新增任务添加到队列首位）
func (q *QueueScheduler) AddHead(id string, execute Execute) {
	logger := log.WithField("add_task_id", id).
		WithField("add_position", "queue_head")
	if id == "" || execute == nil {
		logger.Error("tasks id is empty or tasks execute is nil")
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.cover(id, execute) {
		logger.Error("queue add head failed: tasks already exists")
		return
	}

	newTask := NewQueue(id, execute)
	if q.head != nil {
		// 队列已有任务，将新任务插入到头部
		q.head.prev = newTask
		newTask.next = q.head
		q.head = newTask
	} else {
		// 队列为空，新任务既是头也是尾
		q.head = newTask
		q.tail = newTask
	}
	q.tasks[id] = newTask
	logger.Info("queue add head success")
}

// AddAfter 后插队(将新任务添加到目标任务之后)
func (q *QueueScheduler) AddAfter(baseId, id string, execute Execute) {
	logger := log.WithField("add_task_id", id).
		WithField("base_task_id", baseId).
		WithField("add_position", "after")
	if id == "" || execute == nil {
		logger.Error("tasks id is empty or tasks execute is nil")
		return
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.cover(id, execute) {
		logger.Error("queue add after failed: tasks already exists")
		return
	}

	baseTask, ok := q.tasks[baseId]
	if !ok {
		logger.Error("queue add after failed: base tasks not exist")
		return
	}

	newTask := NewQueue(id, execute)
	if baseTask.next != nil { // 目标任务存在且不是队尾，则插入到目标任务之后
		newTask.next = baseTask.next
		newTask.prev = baseTask
		baseTask.next.prev = newTask
		baseTask.next = newTask
	} else { // 插入当前队尾之后
		baseTask.next = newTask
		newTask.prev = baseTask
		q.tail = newTask
	}
	q.tasks[id] = newTask
	logger.Info("queue add after success")
}

// AddBefore 前插队(将新任务添加到目标任务之后)
func (q *QueueScheduler) AddBefore(baseId, id string, execute Execute) {
	logger := log.WithField("add_task_id", id).
		WithField("base_task_id", baseId).
		WithField("add_position", "before")
	if id == "" || execute == nil {
		logger.Error("tasks id is empty or tasks execute is nil")
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.cover(id, execute) {
		logger.Error("queue add before failed: tasks already exists")
		return
	}

	baseTask, ok := q.tasks[baseId]
	if !ok {
		logger.Error("queue add before failed: base tasks not exist")
		return
	}

	newTask := NewQueue(id, execute)
	if baseTask.prev != nil { // 目标任务不是队列头部，插入到目标任务之前
		newTask.prev = baseTask.prev
		newTask.next = baseTask
		baseTask.prev.next = newTask
		baseTask.prev = newTask
	} else { // 目标任务是队列头部，新任务成为新的头部
		newTask.next = baseTask
		baseTask.prev = newTask
		q.head = newTask
	}
	q.tasks[id] = newTask
	logger.Info("queue add before success")
}

// Remove 移除任务
func (q *QueueScheduler) Remove(id string) {
	if id != "" {
		logger := log.WithField("remove_task_id", id)
		if task, ok := q.tasks[id]; ok {
			q.mutex.Lock()
			defer q.mutex.Unlock()
			if task.prev == nil && task.next == nil {
				// 移除的是队列中唯一的任务
				q.head = nil
				q.tail = nil
			} else if task.prev == nil {
				// 移除头部任务
				q.head = task.next
				if q.head != nil {
					q.head.prev = nil
				}
			} else if task.next == nil {
				// 移除尾部任务
				q.tail = task.prev
				if q.tail != nil {
					q.tail.next = nil
				}
			} else {
				// 移除中间任务
				task.prev.next = task.next
				task.next.prev = task.prev
			}
			delete(q.tasks, id)
			logger.Info("queue remove success")
		} else {
			logger.Error("queue remove failed: execute id does not exist")
		}
	}
}

// 覆盖任务
func (q *QueueScheduler) cover(id string, execute Execute) bool {
	if exist, ok := q.tasks[id]; ok {
		exist.execute = execute
		return true
	}
	return false
}

// Valid 队列是否有效
func (q *QueueScheduler) Valid() bool {
	return q.head != nil && q.tail != nil && len(q.tasks) > 0
}

// Reset 队列重置
func (q *QueueScheduler) Reset() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.head = nil
	q.tail = nil
	q.tasks = make(map[string]*Queue)
}

// Exist 任务是否存在
func (q *QueueScheduler) Exist(id string) bool {
	if _, ok := q.tasks[id]; ok {
		return true
	}
	return false
}

func (q *QueueScheduler) GetIds() []string {
	var ids = make([]string, 0, len(q.tasks))
	for t := q.head; t != nil; t = t.next {
		ids = append(ids, t.id)
	}
	return ids
}
