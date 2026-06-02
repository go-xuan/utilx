package funcx

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// Wrap 包装函数
type Wrap[F any] func(F) F

// WrapV 包装 V 函数
func WrapV(function V, wraps ...Wrap[V]) V {
	if function == nil || len(wraps) == 0 {
		return function
	}
	for _, wrap := range wraps {
		function = wrap(function)
	}
	return function
}

// WrapE 包装 E 函数
func WrapE(function E, wraps ...Wrap[E]) E {
	if function == nil || len(wraps) == 0 {
		return function
	}
	for _, wrap := range wraps {
		function = wrap(function)
	}
	return function
}

// WrapC 包装 C 函数
func WrapC(function C, wraps ...Wrap[C]) C {
	if function == nil || len(wraps) == 0 {
		return function
	}
	for _, wrap := range wraps {
		function = wrap(function)
	}
	return function
}

// WrapX 包装 X 函数
func WrapX(function X, wraps ...Wrap[X]) X {
	if function == nil || len(wraps) == 0 {
		return function
	}
	for _, wrap := range wraps {
		function = wrap(function)
	}
	return function
}

// DurationV 记录 V 函数执行时间
func DurationV(function V) V {
	return func() {
		start := time.Now()
		defer log.WithField("duration", time.Since(start)).Info()
		function()
	}
}

// DurationE 记录 E 函数执行时间
func DurationE(function E) E {
	return func() error {
		start := time.Now()
		defer log.WithField("duration", time.Since(start)).Info()
		return function()
	}
}

// DurationC 记录 C 函数执行时间
func DurationC(function C) C {
	return func(ctx context.Context) {
		start := time.Now()
		defer log.WithField("duration", time.Since(start)).Info()
		function(ctx)
	}
}

// DurationX 记录 X 函数执行时间
func DurationX(function X) X {
	return func(ctx context.Context) error {
		start := time.Now()
		defer log.WithField("duration", time.Since(start)).Info()
		return function(ctx)
	}
}

// RecoverV 捕获 V 函数的 panic
func RecoverV(function V) V {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				log.WithField("panic", r).Error()
			}
		}()
		function()
	}
}

// RecoverE 捕获 E 函数的 panic 并转换为 error
func RecoverE(function E) E {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = logPanic(r)
			}
		}()
		return function()
	}
}

// RecoverX 捕获 X 函数的 panic 并转换为 error
func RecoverX(function X) X {
	return func(ctx context.Context) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = logPanic(r)
			}
		}()
		return function(ctx)
	}
}

func logPanic(r any) error {
	log.WithField("panic", r).Error()
	if err, ok := r.(error); ok {
		return err
	}
	return nil
}
