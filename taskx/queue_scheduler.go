package taskx

import (
	"context"
	"sync"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

// NewQueueScheduler 队列任务处理调度器
func NewQueueScheduler() *QueueScheduler {
	return &QueueScheduler{
		mutex: new(sync.Mutex),
		tasks: make(map[string]*QueueTask),
	}
}

// QueueScheduler 队列任务处理调度器
type QueueScheduler struct {
	mutex *sync.Mutex           // 锁
	head  *QueueTask            // 头部任务
	tail  *QueueTask            // 尾部任务
	tasks map[string]*QueueTask // 任务列表
}

func (q *QueueScheduler) Execute(ctx context.Context) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// 检查队列是否为空
	if q.head == nil {
		log.Info("queue is empty, no tasks to execute")
		return nil
	}

	current := q.head
	for current != nil {
		logger := log.WithField("curr_task_name", current.name)
		if current.next != nil {
			logger = logger.WithField("next_task_name", current.next.name)
		}
		// 执行当前任务
		logger.Info("queue task task")
		if err := current.Execute(ctx); err != nil {
			logger.WithField("error", err.Error()).Error("queue task task error")
			return errorx.Wrap(err, "queue task task error")
		}

		// 从任务列表中删除当前任务并更新当前任务指针
		delete(q.tasks, current.name)
		current = current.next
	}
	return nil
}

// Add 新增队列任务（默认尾插）
func (q *QueueScheduler) Add(name string, task Task) {
	q.AddTail(name, task)
}

// AddTail 尾插（当前新增任务添加到队列末尾）
func (q *QueueScheduler) AddTail(name string, task Task) {
	logger := log.WithField("add_task_name", name).
		WithField("position", "tail")
	if name == "" || task == nil {
		logger.Error("task name is empty or task is nil")
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.cover(name, task) {
		return
	}

	newTask := &QueueTask{name: name, task: task}
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

	q.tasks[name] = newTask
	logger.Info("task add success")
}

// AddHead 头插（当前新增任务添加到队列首位）
func (q *QueueScheduler) AddHead(name string, task Task) {
	logger := log.WithField("add_task_name", name).
		WithField("position", "head")
	if name == "" || task == nil {
		logger.Error("task name is empty or task is nil")
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.cover(name, task) {
		return
	}

	newTask := &QueueTask{name: name, task: task}
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
	q.tasks[name] = newTask
	logger.Info("task add success")
}

// AddAfter 后插队(将新任务添加到目标任务之后)
func (q *QueueScheduler) AddAfter(base, name string, task Task) {
	logger := log.WithField("add_task_name", name).
		WithField("base_task_name", base).
		WithField("position", "after")
	if name == "" || task == nil {
		logger.Error("task name is empty or task is nil")
		return
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.cover(name, task) {
		return
	}

	newTask := &QueueTask{name: name, task: task}
	if curr, ok := q.tasks[base]; ok {
		if curr.next != nil {
			// 目标任务存在且不是队尾，则插入到目标任务之后
			newTask.next = curr.next
			newTask.prev = curr
			curr.next.prev = newTask
			curr.next = newTask
		} else {
			// 插入当前队尾之后
			curr.next = newTask
			newTask.prev = curr
			q.tail = newTask
		}
		q.tasks[name] = newTask
		logger.Info("task add success")
	} else {
		logger.Error("task add failed: base task not exist")
	}
}

// AddBefore 前插队(将新任务添加到目标任务之后)
func (q *QueueScheduler) AddBefore(base, name string, task Task) {
	logger := log.WithField("add_task_name", name).
		WithField("base_task_name", base).
		WithField("position", "before")
	if name == "" || task == nil {
		logger.Error("task name or task is empty")
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.cover(name, task) {
		return
	}

	newTask := &QueueTask{name: name, task: task}
	if curr, ok := q.tasks[base]; ok {
		if curr.prev != nil {
			// 目标任务不是队列头部，插入到目标任务之前
			newTask.prev = curr.prev
			newTask.next = curr
			curr.prev.next = newTask
			curr.prev = newTask
		} else {
			// 目标任务是队列头部，新任务成为新的头部
			newTask.next = curr
			curr.prev = newTask
			q.head = newTask
		}
		q.tasks[name] = newTask
		logger.Info("task add success")
	} else {
		logger.Error("task add failed: base task not exist")
	}
}

// Remove 移除任务
func (q *QueueScheduler) Remove(name string) {
	if name != "" {
		logger := log.WithField("add_task_name", name)
		if task, ok := q.tasks[name]; ok {
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
			delete(q.tasks, name)
			logger.Info("task remove success")
		} else {
			logger.Error("task remove failed: task name does not exist")
		}
	}
}

// 覆盖任务
func (q *QueueScheduler) cover(name string, task Task) bool {
	if existTask, exist := q.tasks[name]; exist {
		log.WithField("add_task_name", name).Errorf("task already exists")
		existTask.task = task
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
	q.tasks = make(map[string]*QueueTask)
}

// Exist 任务是否存在
func (q *QueueScheduler) Exist(name string) bool {
	if _, ok := q.tasks[name]; ok {
		return true
	}
	return false
}

func (q *QueueScheduler) Names() []string {
	var names = make([]string, 0, len(q.tasks))
	for t := q.head; t != nil; t = t.next {
		names = append(names, t.name)
	}
	return names
}
