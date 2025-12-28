package funcx

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// Wrap 包装函数
type Wrap[F any] func(F) F

// WrapV 包装V函数
func WrapV(function V, wraps ...Wrap[V]) V {
	if function != nil && len(wraps) > 0 {
		for _, wrap := range wraps {
			function = wrap(function)
		}
	}
	return function
}

// WrapE 包装E函数
func WrapE(function E, wraps ...Wrap[E]) E {
	if function != nil && len(wraps) > 0 {
		for _, wrap := range wraps {
			function = wrap(function)
		}
	}
	return function
}

// WrapC 包装C函数
func WrapC(function C, wraps ...Wrap[C]) C {
	if function != nil && len(wraps) > 0 {
		for _, wrap := range wraps {
			function = wrap(function)
		}
	}
	return function
}

// WrapX 包装X函数
func WrapX(function X, wraps ...Wrap[X]) X {
	if function != nil && len(wraps) > 0 {
		for _, wrap := range wraps {
			function = wrap(function)
		}
	}
	return function
}

// DurationV 记录V函数执行时间
func DurationV(function V) V {
	return func() {
		start := time.Now()
		defer log.WithField("duration", time.Since(start)).Info()
		function()
	}
}

// DurationE 记录E函数执行时间
func DurationE(function E) E {
	return func() error {
		start := time.Now()
		defer log.WithField("duration", time.Since(start)).Info()
		return function()
	}
}

// DurationC 记录C函数执行时间
func DurationC(function C) C {
	return func(ctx context.Context) {
		start := time.Now()
		defer log.WithField("duration", time.Since(start)).Info()
		function(ctx)
	}
}

// DurationX 记录X函数执行时间
func DurationX(function X) X {
	return func(ctx context.Context) error {
		start := time.Now()
		defer log.WithField("duration", time.Since(start)).Info()
		return function(ctx)
	}
}
