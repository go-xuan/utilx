package idx

import (
	"sync"

	"github.com/go-xuan/typex"
)

var sequencePool *SequencePool

// GetSequencePool 获取序列池
func GetSequencePool() *SequencePool {
	if sequencePool == nil {
		sequencePool = &SequencePool{
			Pool: typex.NewStringEnum[*Sequence](),
		}
	}
	return sequencePool
}

// SequencePool 序列池
type SequencePool struct {
	Pool *typex.Enum[string, *Sequence]
}

// Create 创建序列
func (m *SequencePool) Create(name string, start, incr int64) {
	m.Pool.Add(name, NewSequence(name, start, incr))
}

// CurrVal 获取序列当前值
func (m *SequencePool) CurrVal(name string) int64 {
	if sequence := m.Pool.Get(name); sequence != nil {
		return sequence.Curr()
	}
	m.Create(name, 0, 1)
	return 0
}

// NextVal 获取序列下一个值
func (m *SequencePool) NextVal(name string) int64 {
	if sequence := m.Pool.Get(name); sequence != nil {
		return sequence.Next()
	}
	m.Create(name, 1, 1)
	return 1
}

// NextBatch 获取序列当前值
func (m *SequencePool) NextBatch(name string, size int64) int64 {
	if sequence := m.Pool.Get(name); sequence != nil {
		var next = sequence.Next()
		sequence.Set(next + (size-1)*sequence.incr)
		return next
	}
	m.Create(name, size+1, 1)
	return 1
}

// Set 设置序列当前值
func (m *SequencePool) Set(name string, value int64) {
	if sequence := m.Pool.Get(name); sequence != nil {
		sequence.Set(value)
		return
	}
	m.Create(name, value, 1)
}

// Reset 序列重置
func (m *SequencePool) Reset(name string) {
	if sequence := m.Pool.Get(name); sequence != nil {
		sequence.Reset()
		return
	}
	m.Create(name, 0, 1)
}

// NewSequence 新建序列
func NewSequence(name string, start, incr int64) *Sequence {
	return &Sequence{
		sync.RWMutex{},
		name,
		start,
		incr,
		start,
	}
}

// Sequence 序列
type Sequence struct {
	sync.RWMutex        // 读写锁
	name         string // 序列名
	start        int64  // 开始值
	incr         int64  // 递增值
	val          int64  // 序列号
}

// Next 获取序列值
func (s *Sequence) Next() int64 {
	s.Lock()
	defer s.Unlock()
	s.val += s.incr
	return s.val
}

// Curr 获取序列当前值
func (s *Sequence) Curr() int64 {
	s.RLock()
	defer s.RUnlock()
	return s.val
}

// Set 设置序列当前值
func (s *Sequence) Set(v int64) {
	s.Lock()
	defer s.Unlock()
	if v < s.start {
		s.val = s.start
		return
	}
	s.val = v
}

// Reset 序列重置
func (s *Sequence) Reset() {
	s.Lock()
	defer s.Unlock()
	s.val = s.start
}
