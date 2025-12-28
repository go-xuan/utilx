package errorx

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

// Close 关闭closer并收集错误
func Close(closer io.Closer, collector ...Collector) {
	if closer != nil {
		Collect(closer.Close(), collector...)
	}
}

// Collect 收集错误
func Collect(err error, collector ...Collector) {
	if len(collector) > 0 {
		for _, c := range collector {
			c.Collect(err)
		}
	}
}

// Collector 错误收集器
type Collector interface {
	Collect(err error)
}

// NewPrintCollector 创建打印收集器
func NewPrintCollector() Collector {
	return Print{}
}

// Print 打印收集器
type Print struct{}

func (p Print) Collect(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// NewLogCollector 创建日志收集器
func NewLogCollector() Collector {
	return &Log{
		Logger: log.StandardLogger(),
		Level:  log.ErrorLevel,
	}
}

// Log 日志收集器
type Log struct {
	Logger *log.Logger
	Level  log.Level
}

func (l *Log) Collect(err error) {
	if err != nil {
		l.setDefault()
		l.Logger.Log(l.Level, err)
	}
}

// 设置默认日志记录器和日志级别
func (l *Log) setDefault() {
	if l.Logger == nil {
		l.Logger = log.StandardLogger()
	}
	if l.Level == 0 {
		l.Level = log.ErrorLevel
	}
}
