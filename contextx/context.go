package contextx

import (
	"context"

	"github.com/go-xuan/typex"
)

const (
	valuesKey = "__x_values__"
)

// New 新建上下文
func New() context.Context {
	var values = make(map[string]any)
	return context.WithValue(context.Background(), valuesKey, values)
}

// SetValue 设置值
func SetValue(ctx context.Context, key string, value any) {
	if v := ctx.Value(valuesKey); v != nil {
		if values, ok := v.(map[string]any); ok {
			values[key] = value
		}
	}
}

// GetValue 获取值
func GetValue(ctx context.Context, key string) typex.Value {
	if value := getValue(ctx, key); value != nil {
		return typex.NewValue(value)
	}
	return typex.NewZero()
}

func getValue(ctx context.Context, key string) any {
	if v := ctx.Value(valuesKey); v != nil {
		if values, ok := v.(map[string]any); ok {
			return values[key]
		}
	}
	return nil
}
