package taskx

import (
	"context"

	"github.com/go-xuan/utilx/errorx"
)

// NewSerial 新建串行器
func NewSerial(total, limit int, execute func(context.Context, int, int) error) *Serial {
	return &Serial{
		total:   total,
		limit:   limit,
		execute: execute,
	}
}

// Serial 串行器
type Serial struct {
	total   int
	limit   int
	offset  int
	execute func(context.Context, int, int) error
}

func (c *Serial) Execute(ctx context.Context) error {
	for c.offset < c.total {
		next := c.offset + c.limit
		if next > c.total {
			next = c.total
		}
		if err := c.execute(ctx, c.offset, next); err != nil {
			return errorx.Wrap(err, "cursor execute error")
		}
		c.offset = next
	}
	return nil
}
