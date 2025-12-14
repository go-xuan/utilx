package taskx

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-xuan/utilx/funcx"
	"github.com/robfig/cron/v3"
)

func TestCronScheduler(t *testing.T) {
	logger := noLogger{}
	scheduler := NewCronScheduler("cron_scheduler_test",
		cron.WithLogger(logger),
		cron.WithChain(cron.SkipIfStillRunning(logger)),
	)

	scheduler.AddWrap(funcx.XDuration)
	scheduler.AddJob("1", "@every 5s", function(1, 0.5))
	scheduler.AddJob("2", "@every 2s", function(2, 0.8))
	scheduler.AddJob("3", "@daily", function(3, 0.3))
	scheduler.AddJob("4", "0 */1 * * * ?", function(4, 0.7))

	// 开始调度
	if err := scheduler.Start(); err != nil {
		t.Log(err)
		return
	}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		for _, job := range scheduler.All() {
			fmt.Println(job.GetMeta())
		}
	}
	select {}
}

type noLogger struct{}

func (n noLogger) Info(string, ...interface{}) {}

func (n noLogger) Error(error, string, ...interface{}) {}
