package taskx

import (
	"context"

	"github.com/go-xuan/utilx/errorx"
)

// Wrap 任务包装
type Wrap func(Execute) Execute

// Wrapper 包装器
type Wrapper struct {
	wraps []Wrap
}

// Add 添加包装器
func (w *Wrapper) Add(wraps ...Wrap) {
	if len(wraps) > 0 {
		w.wraps = append(w.wraps, wraps...)
	}
}

// Wrap 包装函数
func (w *Wrapper) Wrap(execute Execute) Execute {
	if len(w.wraps) > 0 {
		for _, wrap := range w.wraps {
			execute = wrap(execute)
		}
	}
	return execute
}

// Execute 执行
func (w *Wrapper) Execute(ctx context.Context, execute Execute) error {
	if err := w.Wrap(execute)(ctx); err != nil {
		return errorx.Wrap(err, "wrapper execute error")
	}
	return nil
}
