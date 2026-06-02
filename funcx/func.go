// Package funcx 提供函数包装和执行工具
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

// ExecuteV 执行多个 V 函数
func ExecuteV(functions ...V) {
	for _, function := range functions {
		function()
	}
}

// ExecuteE 依次执行多个 E 函数，任意某个执行错误则返回
func ExecuteE(functions ...E) error {
	for _, function := range functions {
		if err := function(); err != nil {
			return errorx.Wrap(err, "execute failed")
		}
	}
	return nil
}

// ExecuteC 执行多个 C 函数
func ExecuteC(ctx context.Context, functions ...C) {
	for _, function := range functions {
		function(ctx)
	}
}

// ExecuteX 依次执行多个 X 函数，任意某个函数执行错误则返回
func ExecuteX(ctx context.Context, functions ...X) error {
	for _, function := range functions {
		if err := function(ctx); err != nil {
			return errorx.Wrap(err, "execute failed")
		}
	}
	return nil
}

// LogErrorE 依次执行多个 E 函数，任意某个函数执行错误则记录日志
func LogErrorE(functions ...E) {
	for _, function := range functions {
		if err := function(); err != nil {
			log.WithError(err).Error()
		}
	}
}

// LogErrorX 执行多个 X 函数，任意某个函数执行错误则记录日志
func LogErrorX(ctx context.Context, functions ...X) {
	for _, function := range functions {
		if err := function(ctx); err != nil {
			log.WithError(err).Error()
		}
	}
}

// PanicE 执行多个 E 函数，任意某个函数执行错误则 panic
func PanicE(functions ...E) {
	for _, function := range functions {
		if err := function(); err != nil {
			panic(err)
		}
	}
}

// PanicX 执行多个 X 函数，任意某个函数执行错误则 panic
func PanicX(ctx context.Context, functions ...X) {
	for _, function := range functions {
		if err := function(ctx); err != nil {
			panic(err)
		}
	}
}

// MergeV 合并多个 V 函数
func MergeV(functions ...V) V {
	return func() {
		for _, function := range functions {
			function()
		}
	}
}

// MergeE 合并多个 E 函数
func MergeE(functions ...E) E {
	return func() error {
		for _, function := range functions {
			if err := function(); err != nil {
				return errorx.Wrap(err, "execute failed")
			}
		}
		return nil
	}
}

// MergeC 合并多个 C 函数
func MergeC(functions ...C) C {
	return func(ctx context.Context) {
		for _, function := range functions {
			function(ctx)
		}
	}
}

// MergeX 合并多个 X 函数
func MergeX(functions ...X) X {
	return func(ctx context.Context) error {
		for _, function := range functions {
			if err := function(ctx); err != nil {
				return errorx.Wrap(err, "execute failed")
			}
		}
		return nil
	}
}
