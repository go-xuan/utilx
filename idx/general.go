package idx

import (
	"time"

	"github.com/google/uuid"
)

// UUID 生成UUID
func UUID() string {
	return uuid.NewString()
}

// TimeUnix 生成时间戳(毫秒)
func TimeUnix() int64 {
	return time.Now().UnixMilli()
}

// Timestamp 生成时间戳
func Timestamp() string {
	return time.Now().Format("20060102150405")
}
