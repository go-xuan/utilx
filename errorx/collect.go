package errorx

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

// Close closes the closer and collects the error with collector.
func Close(closer io.Closer, collector ...Collector) {
	if closer != nil {
		Collect(closer.Close(), collector...)
	}
}

// Collect collect errors with collector.
func Collect(err error, collector ...Collector) {
	if len(collector) > 0 {
		for _, c := range collector {
			c.Collect(err)
		}
	}
}

type Collector interface {
	Collect(err error)
}

func NewPrintCollector() Collector {
	return Print{}
}

type Print struct{}

func (p Print) Collect(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// NewLogCollector returns a collector that logs errors at ErrorLevel.
func NewLogCollector() Collector {
	return &Log{
		Logger: log.StandardLogger(),
		Level:  log.ErrorLevel,
	}
}

type Log struct {
	Logger *log.Logger
	Level  log.Level
}

func (l *Log) setDefault() {
	if l.Logger == nil {
		l.Logger = log.StandardLogger()
	}
	if l.Level == 0 {
		l.Level = log.ErrorLevel
	}
}

func (l *Log) Collect(err error) {
	if err != nil {
		l.setDefault()
		l.Logger.Log(l.Level, err)
	}
}
