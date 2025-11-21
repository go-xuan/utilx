package funcx

import (
	"context"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

type (
	V func()                      // 无参数无返回值函数
	E func() error                // 无参数返回错误的函数
	C func(context.Context)       // 有上下文参数无返回值函数
	X func(context.Context) error // 有上下文参数返回错误的函数
)

type (
	WarpV func(V) V // 无参数无返回值函数包装器
	WarpE func(E) E // 无参数返回错误的函数包装器
	WarpC func(C) C // 有上下文参数无返回值函数包装器
	WarpX func(X) X // 有上下文参数返回错误的函数包装器
)

// ExecuteV 执行多个函数
func ExecuteV(fns ...V) {
	for _, fn := range fns {
		fn()
	}
}

// ExecuteE 依次执行多个函数，任意某个执行错误则返回
func ExecuteE(fns ...E) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return errorx.Wrap(err, "execute failed")
		}
	}
	return nil
}

// ExecuteC 执行多个函数
func ExecuteC(ctx context.Context, fns ...C) {
	for _, fn := range fns {
		fn(ctx)
	}
}

// ExecuteX 依次执行多个函数，任意某个函数执行错误则返回
func ExecuteX(ctx context.Context, fns ...X) error {
	for _, fn := range fns {
		if err := fn(ctx); err != nil {
			return errorx.Wrap(err, "execute failed")
		}
	}
	return nil
}

// LogE 依次执行多个函数，任意某个函数执行错误则记录日志
func LogE(fns ...E) {
	for _, fn := range fns {
		if err := fn(); err != nil {
			log.WithError(err).Error()
		}
	}
}

// LogX 执行多个函数，任意某个函数执行错误则记录日志
func LogX(ctx context.Context, fns ...X) {
	for _, fn := range fns {
		if err := fn(ctx); err != nil {
			log.WithError(err).Error()
		}
	}
}

// PanicE 执行多个函数, 任意某个函数执行错误则panic
func PanicE(fns ...E) {
	for _, fn := range fns {
		if err := fn(); err != nil {
			panic(err)
		}
	}
}

// PanicX 执行多个函数, 任意某个函数执行错误则panic
func PanicX(ctx context.Context, fns ...X) {
	for _, fn := range fns {
		if err := fn(ctx); err != nil {
			panic(err)
		}
	}
}

// MergeV 合并多个函数
func MergeV(fns ...V) V {
	return func() {
		for _, fn := range fns {
			fn()
		}
	}
}

// MergeE 合并多个函数
func MergeE(fns ...E) E {
	return func() error {
		for _, fn := range fns {
			if err := fn(); err != nil {
				return errorx.Wrap(err, "execute failed")
			}
		}
		return nil
	}
}

// MergeC 合并多个函数
func MergeC(fns ...C) C {
	return func(ctx context.Context) {
		for _, fn := range fns {
			fn(ctx)
		}
	}
}

// MergeX 合并多个函数
func MergeX(ctx context.Context, fns ...X) X {
	return func(context.Context) error {
		for _, fn := range fns {
			if err := fn(ctx); err != nil {
				return errorx.Wrap(err, "execute failed")
			}
		}
		return nil
	}
}
