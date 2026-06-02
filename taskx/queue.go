package taskx

import (
	"context"
	"sync"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/funcx"
	log "github.com/sirupsen/logrus"
)

// NewQueue 队列任务处理调度器
func NewQueue(name string) *Queue {
	return &Queue{
		name:  name,
		mutex: new(sync.Mutex),
		tasks: make(map[string]*queueTask),
	}
}

// queueTask 队列任务
type queueTask struct {
	name     string     // 任务ID
	function funcx.X    // 任务执行函数
	prev     *queueTask // 指向上一个任务
	next     *queueTask // 指向下一个任务
}

// Queue 队列任务处理调度器
type Queue struct {
	name  string                // 队列任务处理调度器ID
	mutex *sync.Mutex           // 锁
	head  *queueTask            // 头部任务
	tail  *queueTask            // 尾部任务
	tasks map[string]*queueTask // 任务列表
}

func (q *Queue) GetID() string {
	return q.name
}

func (q *Queue) Execute(ctx context.Context) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	logger := log.WithField("queue", q.GetID())
	// 检查队列是否为空
	if q.head == nil || len(q.tasks) == 0 {
		logger.Info("queue tasks is empty")
		return nil
	}

	// 遍历执行队列中的任务
	current := q.head
	for current != nil {
		logger = logger.WithField("current", current.name)
		if current.next != nil {
			logger = logger.WithField("next", current.next.name)
		}

		if err := current.function(ctx); err != nil {
			logger.WithError(err).Error("queue task execute error")
			return errorx.Wrap(err, "queue task execute error")
		}
		logger.Info("queue task func execute success")
		// 从任务列表中删除当前任务并更新当前任务指针
		delete(q.tasks, current.name)
		current = current.next
	}
	return nil
}

// Add 新增任务（默认尾插）
func (q *Queue) Add(name string, fn funcx.X) {
	q.AddTail(name, fn)
}

// AddTail 尾插（当前新增任务添加到队列末尾）
func (q *Queue) AddTail(name string, fn funcx.X) {
	logger := log.WithField("name", name).
		WithField("queue", q.GetID()).
		WithField("position", "tail of queue")

	if name == "" || fn == nil {
		logger.Error("tasks name is empty or tasks func is nil")
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.cover(name, fn) {
		logger.Error("queue add tail failed, tasks already exists")
		return
	}

	task := &queueTask{name: name, function: fn}
	if q.tail != nil {
		// 队列不为空，将新任务添加到尾部
		task.prev = q.tail
		q.tail.next = task
		q.tail = task
	} else {
		// 队列为空，新任务既是头也是尾
		q.head = task
		q.tail = task
	}

	q.tasks[name] = task
	logger.Info("queue add tail success")
}

// AddHead 头插（当前新增任务添加到队列首位）
func (q *Queue) AddHead(name string, fn funcx.X) {
	logger := log.WithField("name", name).
		WithField("queue", q.GetID()).
		WithField("position", "head of queue")

	if name == "" || fn == nil {
		logger.Error("tasks name is empty or tasks func is nil")
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.cover(name, fn) {
		logger.Error("queue add head failed: tasks already exists")
		return
	}

	newTask := &queueTask{name: name, function: fn}
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
	logger.Info("queue add head success")
}

// AddAfter 后插队(将新任务添加到目标任务之后)
func (q *Queue) AddAfter(baseName, name string, fn funcx.X) {
	logger := log.WithField("name", name).
		WithField("queue", q.GetID()).
		WithField("position", "after task: "+baseName)
	if name == "" || fn == nil {
		logger.Error("task name is empty or task func is nil")
		return
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.cover(name, fn) {
		logger.Error("queue add after failed: tasks already exists")
		return
	}

	baseTask, ok := q.tasks[baseName]
	if !ok {
		logger.Error("queue add after failed: base tasks not exist")
		return
	}

	newTask := &queueTask{name: name, function: fn}
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
	q.tasks[name] = newTask
	logger.Info("queue add after success")
}

// AddBefore 前插队(将新任务添加到目标任务之后)
func (q *Queue) AddBefore(baseName, name string, fn funcx.X) {
	logger := log.WithField("name", name).
		WithField("queue", q.GetID()).
		WithField("position", "before task: "+baseName)
	if name == "" || fn == nil {
		logger.Error("tasks name is empty or tasks func is nil")
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.cover(name, fn) {
		logger.Error("queue add before failed: task already exists")
		return
	}

	baseTask, ok := q.tasks[baseName]
	if !ok {
		logger.Error("queue add before failed: base tasks not exist")
		return
	}

	newTask := &queueTask{name: name, function: fn}
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
	q.tasks[name] = newTask
	logger.Info("queue add before success")
}

// Remove 移除任务
func (q *Queue) Remove(name string) {
	if name != "" {
		logger := log.WithField("name", name)
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if task, ok := q.tasks[name]; ok {
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
			logger.Info("queue remove success")
		} else {
			logger.Error("queue remove failed: task name does not exist")
		}
	}
}

// 覆盖任务
func (q *Queue) cover(name string, fn funcx.X) bool {
	if exist, ok := q.tasks[name]; ok {
		exist.function = fn
		return true
	}
	return false
}

// Valid 队列是否有效
func (q *Queue) Valid() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.head != nil && q.tail != nil && len(q.tasks) > 0
}

// Reset 队列重置
func (q *Queue) Reset() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.head = nil
	q.tail = nil
	q.tasks = make(map[string]*queueTask)
}

// Exist 任务是否存在
func (q *Queue) Exist(name string) bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	_, ok := q.tasks[name]
	return ok
}

func (q *Queue) GetNames() []string {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	names := make([]string, 0, len(q.tasks))
	for t := q.head; t != nil; t = t.next {
		names = append(names, t.name)
	}
	return names
}
