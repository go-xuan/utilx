package bytex

import (
	"fmt"
)

const (
	KB = 1 << 10 // 1024
	MB = 1 << 20 // 1048576
	GB = 1 << 30 // 1073741824
	TB = 1 << 40
)

// Convert 自动适配单位
func Convert(bytes int64, decimal int) string {
	switch {
	case bytes >= TB:
		return Byte2TB(bytes, decimal)
	case bytes >= GB:
		return Byte2GB(bytes, decimal)
	case bytes >= MB:
		return Byte2MB(bytes, decimal)
	case bytes >= KB:
		return Byte2KB(bytes, decimal)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// Byte2KB 强制转换字节为KB
func Byte2KB(bytes int64, decimal int) string {
	formatStr := fmt.Sprintf("%%.%df KB", decimal)
	return fmt.Sprintf(formatStr, float64(bytes)/MB)
}

// Byte2MB 强制转换字节为MB
func Byte2MB(bytes int64, decimal int) string {
	formatStr := fmt.Sprintf("%%.%df MB", decimal)
	return fmt.Sprintf(formatStr, float64(bytes)/MB)
}

// Byte2GB 强制转换字节为GB
func Byte2GB(bytes int64, decimal int) string {
	formatStr := fmt.Sprintf("%%.%df GB", decimal)
	return fmt.Sprintf(formatStr, float64(bytes)/GB)
}

// Byte2TB 强制转换字节为TB
func Byte2TB(bytes int64, decimal int) string {
	formatStr := fmt.Sprintf("%%.%df TB", decimal)
	return fmt.Sprintf(formatStr, float64(bytes)/TB)
}
