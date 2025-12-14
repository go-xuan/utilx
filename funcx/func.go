package funcx

import (
	"context"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

type (
	V func()                      // 无参数无返回值函数
	C func(context.Context)       // 有上下文参数无返回值函数
	E func() error                // 无参数返回错误的函数
	X func(context.Context) error // 有上下文参数且返回错误的函数
)

// VExecute 执行多个V函数
func VExecute(functions ...V) {
	for _, function := range functions {
		function()
	}
}

// EExecute 依次执行多个E函数，任意某个执行错误则返回
func EExecute(functions ...E) error {
	for _, function := range functions {
		if err := function(); err != nil {
			return errorx.Wrap(err, "execute failed")
		}
	}
	return nil
}

// CExecute 执行多个C函数
func CExecute(ctx context.Context, functions ...C) {
	for _, function := range functions {
		function(ctx)
	}
}

// XExecute 依次执行多个X函数，任意某个函数执行错误则返回
func XExecute(ctx context.Context, functions ...X) error {
	for _, function := range functions {
		if err := function(ctx); err != nil {
			return errorx.Wrap(err, "execute failed")
		}
	}
	return nil
}

// EErrorLog 依次执行多个E函数，任意某个函数执行错误则记录日志
func EErrorLog(functions ...E) {
	for _, function := range functions {
		if err := function(); err != nil {
			log.WithError(err).Error()
		}
	}
}

// XErrorLog 执行多个X函数，任意某个函数执行错误则记录日志
func XErrorLog(ctx context.Context, functions ...X) {
	for _, function := range functions {
		if err := function(ctx); err != nil {
			log.WithError(err).Error()
		}
	}
}

// EPanic 执行多个E函数, 任意某个函数执行错误则panic
func EPanic(functions ...E) {
	for _, function := range functions {
		if err := function(); err != nil {
			panic(err)
		}
	}
}

// XPanic 执行多个X函数, 任意某个函数执行错误则panic
func XPanic(ctx context.Context, functions ...X) {
	for _, function := range functions {
		if err := function(ctx); err != nil {
			panic(err)
		}
	}
}

// VMerge 合并多个V函数
func VMerge(functions ...V) V {
	return func() {
		for _, function := range functions {
			function()
		}
	}
}

// EMerge 合并多个E函数
func EMerge(functions ...E) E {
	return func() error {
		for _, function := range functions {
			if err := function(); err != nil {
				return errorx.Wrap(err, "execute failed")
			}
		}
		return nil
	}
}

// CMerge 合并多个C函数
func CMerge(functions ...C) C {
	return func(ctx context.Context) {
		for _, function := range functions {
			function(ctx)
		}
	}
}

// XMerge 合并多个X函数
func XMerge(ctx context.Context, functions ...X) X {
	return func(context.Context) error {
		for _, function := range functions {
			if err := function(ctx); err != nil {
				return errorx.Wrap(err, "execute failed")
			}
		}
		return nil
	}
}
