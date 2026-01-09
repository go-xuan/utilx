package idx

import (
	"sync"

	"github.com/go-xuan/typex"
)

var sequencePool *typex.Enum[string, *Sequence]

// GetSequence 获取序列
func GetSequence(name string) *Sequence {
	if sequencePool == nil {
		sequencePool = typex.NewStringEnum[*Sequence]()
	}
	sequence, ok := sequencePool.Exist(name)
	if !ok {
		sequence = NewSequence(name, 0, 1)
		sequencePool.Add(name, sequence)
	}
	return sequence
}

// AddSequence 创建序列
func AddSequence(name string, start, incr int64) *Sequence {
	if sequencePool == nil {
		sequencePool = typex.NewStringEnum[*Sequence]()
	}
	sequence := NewSequence(name, start, incr)
	sequencePool.Add(name, sequence)
	return sequence
}

// NewSequence 新建序列
func NewSequence(name string, start, incr int64) *Sequence {
	return &Sequence{
		RWMutex: sync.RWMutex{},
		name:    name,
		val:     start,
		start:   start,
		incr:    incr,
	}
}

// Sequence 序列
type Sequence struct {
	sync.RWMutex        // 读写锁
	name         string // 序列名
	val          int64  // 序列号
	start        int64  // 开始值
	incr         int64  // 递增值
}

// Curr 获取序列当前值
func (s *Sequence) Curr() int64 {
	s.RLock()
	defer s.RUnlock()
	return s.val
}

// Next 获取序列值
func (s *Sequence) Next() int64 {
	s.Lock()
	defer s.Unlock()
	s.val += s.incr
	return s.val
}

// Set 设置序列当前值
func (s *Sequence) Set(val int64) {
	s.Lock()
	defer s.Unlock()
	if val < s.start {
		val = s.start
	}
	s.val = val
}

// Reset 序列重置
func (s *Sequence) Reset() {
	s.Lock()
	defer s.Unlock()
	s.val = s.start
}
