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

func (s *QueueScheduler) Execute(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 检查队列是否为空
	if s.head == nil {
		log.Info("queue is empty, no execute to execute")
		return nil
	}

	current := s.head
	for current != nil {
		logger := log.WithField("curr_task_name", current.name)
		if current.next != nil {
			logger = logger.WithField("next_task_name", current.next.name)
		}
		// 执行当前任务
		logger.Info("queue execute execute")
		if err := current.Execute(ctx); err != nil {
			logger.WithField("error", err.Error()).Error("queue execute execute error")
			return errorx.Wrap(err, "queue execute execute error")
		}

		// 从任务列表中删除当前任务并更新当前任务指针
		delete(s.tasks, current.name)
		current = current.next
	}
	return nil
}

// Add 新增任务（默认尾插）
func (s *QueueScheduler) Add(name string, execute Execute) {
	s.AddTail(name, execute)
}

// AddTail 尾插（当前新增任务添加到队列末尾）
func (s *QueueScheduler) AddTail(name string, execute Execute) {
	logger := log.WithField("add_task_name", name).
		WithField("position", "tail")
	if name == "" || execute == nil {
		logger.Error("execute name is empty or execute is nil")
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.cover(name, execute) {
		return
	}

	newTask := &QueueTask{name: name, execute: execute}
	if s.tail != nil {
		// 队列不为空，将新任务添加到尾部
		newTask.prev = s.tail
		s.tail.next = newTask
		s.tail = newTask
	} else {
		// 队列为空，新任务既是头也是尾
		s.head = newTask
		s.tail = newTask
	}

	s.tasks[name] = newTask
	logger.Info("execute add success")
}

// AddHead 头插（当前新增任务添加到队列首位）
func (s *QueueScheduler) AddHead(name string, execute Execute) {
	logger := log.WithField("add_task_name", name).
		WithField("position", "head")
	if name == "" || execute == nil {
		logger.Error("execute name is empty or execute is nil")
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.cover(name, execute) {
		return
	}

	newTask := &QueueTask{name: name, execute: execute}
	if s.head != nil {
		// 队列已有任务，将新任务插入到头部
		s.head.prev = newTask
		newTask.next = s.head
		s.head = newTask
	} else {
		// 队列为空，新任务既是头也是尾
		s.head = newTask
		s.tail = newTask
	}
	s.tasks[name] = newTask
	logger.Info("execute add success")
}

// AddAfter 后插队(将新任务添加到目标任务之后)
func (s *QueueScheduler) AddAfter(base, name string, execute Execute) {
	logger := log.WithField("add_task_name", name).
		WithField("base_task_name", base).
		WithField("position", "after")
	if name == "" || execute == nil {
		logger.Error("execute name is empty or execute is nil")
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.cover(name, execute) {
		return
	}

	newTask := &QueueTask{name: name, execute: execute}
	if curr, ok := s.tasks[base]; ok {
		if curr.next != nil { // 目标任务存在且不是队尾，则插入到目标任务之后
			newTask.next = curr.next
			newTask.prev = curr
			curr.next.prev = newTask
			curr.next = newTask
		} else { // 插入当前队尾之后
			curr.next = newTask
			newTask.prev = curr
			s.tail = newTask
		}
		s.tasks[name] = newTask
		logger.Info("add success")
	} else {
		logger.Error("add failed: base execute not exist")
	}
}

// AddBefore 前插队(将新任务添加到目标任务之后)
func (s *QueueScheduler) AddBefore(base, name string, execute Execute) {
	logger := log.WithField("add_task_name", name).
		WithField("base_task_name", base).
		WithField("position", "before")
	if name == "" || execute == nil {
		logger.Error("execute name or execute is empty")
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.cover(name, execute) {
		return
	}

	newTask := &QueueTask{name: name, execute: execute}
	if curr, ok := s.tasks[base]; ok {
		if curr.prev != nil { // 目标任务不是队列头部，插入到目标任务之前
			newTask.prev = curr.prev
			newTask.next = curr
			curr.prev.next = newTask
			curr.prev = newTask
		} else { // 目标任务是队列头部，新任务成为新的头部
			newTask.next = curr
			curr.prev = newTask
			s.head = newTask
		}
		s.tasks[name] = newTask
		logger.Info("add success")
	} else {
		logger.Error("add failed: base execute not exist")
	}
}

// Remove 移除任务
func (s *QueueScheduler) Remove(name string) {
	if name != "" {
		logger := log.WithField("add_task_name", name)
		if task, ok := s.tasks[name]; ok {
			s.mutex.Lock()
			defer s.mutex.Unlock()
			if task.prev == nil && task.next == nil {
				// 移除的是队列中唯一的任务
				s.head = nil
				s.tail = nil
			} else if task.prev == nil {
				// 移除头部任务
				s.head = task.next
				if s.head != nil {
					s.head.prev = nil
				}
			} else if task.next == nil {
				// 移除尾部任务
				s.tail = task.prev
				if s.tail != nil {
					s.tail.next = nil
				}
			} else {
				// 移除中间任务
				task.prev.next = task.next
				task.next.prev = task.prev
			}
			delete(s.tasks, name)
			logger.Info("remove success")
		} else {
			logger.Error("remove failed: execute name does not exist")
		}
	}
}

// 覆盖任务
func (s *QueueScheduler) cover(name string, execute Execute) bool {
	if exist, ok := s.tasks[name]; ok {
		log.WithField("add_task_name", name).Errorf("execute already exists")
		exist.execute = execute
		return true
	}
	return false
}

// Valid 队列是否有效
func (s *QueueScheduler) Valid() bool {
	return s.head != nil && s.tail != nil && len(s.tasks) > 0
}

// Reset 队列重置
func (s *QueueScheduler) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.head = nil
	s.tail = nil
	s.tasks = make(map[string]*QueueTask)
}

// Exist 任务是否存在
func (s *QueueScheduler) Exist(name string) bool {
	if _, ok := s.tasks[name]; ok {
		return true
	}
	return false
}

func (s *QueueScheduler) Names() []string {
	var names = make([]string, 0, len(s.tasks))
	for t := s.head; t != nil; t = t.next {
		names = append(names, t.name)
	}
	return names
}
