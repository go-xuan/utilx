package taskx

import (
	"context"

	"github.com/go-xuan/utilx/errorx"
)

// NewSplitter 新建拆分策略
func NewSplitter(limit int) *Splitter {
	return &Splitter{
		limit: limit,
	}
}

// Splitter 拆分策略
type Splitter struct {
	limit int // 每批执行任务数
}

func (s *Splitter) Execute(ctx context.Context, total int, execute func(ctx context.Context, start, end, times int) error) error {
	var offset, times int
	for offset < total {
		times++
		next := offset + s.limit
		if next > total {
			next = total
		}
		if err := execute(ctx, offset, next, times); err != nil {
			return errorx.Wrap(err, "splitter execute error")
		}
		offset = next
	}
	return nil
}
