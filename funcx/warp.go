package funcx

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// Warp 包装函数
type Warp[F any] func(F) F

// VWarp 包装V函数
func VWarp(function V, wraps ...Warp[V]) V {
	if function != nil && len(wraps) > 0 {
		for _, wrap := range wraps {
			function = wrap(function)
		}
	}
	return function
}

// EWarp 包装E函数
func EWarp(function E, wraps ...Warp[E]) E {
	if function != nil && len(wraps) > 0 {
		for _, wrap := range wraps {
			function = wrap(function)
		}
	}
	return function
}

// CWarp 包装C函数
func CWarp(function C, wraps ...Warp[C]) C {
	if function != nil && len(wraps) > 0 {
		for _, wrap := range wraps {
			function = wrap(function)
		}
	}
	return function
}

// XWarp 包装X函数
func XWarp(function X, wraps ...Warp[X]) X {
	if function != nil && len(wraps) > 0 {
		for _, wrap := range wraps {
			function = wrap(function)
		}
	}
	return function
}

// VDuration 记录V函数执行时间
func VDuration(function V) V {
	return func() {
		start := time.Now()
		defer log.WithField("duration", time.Since(start)).Info()
		function()
	}
}

// EDuration 记录E函数执行时间
func EDuration(function E) E {
	return func() error {
		start := time.Now()
		defer log.WithField("duration", time.Since(start)).Info()
		return function()
	}
}

// CDuration 记录C函数执行时间
func CDuration(function C) C {
	return func(ctx context.Context) {
		start := time.Now()
		defer log.WithField("duration", time.Since(start)).Info()
		function(ctx)
	}
}

// XDuration 记录X函数执行时间
func XDuration(function X) X {
	return func(ctx context.Context) error {
		start := time.Now()
		defer log.WithField("duration", time.Since(start)).Info()
		return function(ctx)
	}
}
