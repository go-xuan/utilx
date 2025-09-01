package taskx

import (
	"fmt"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	// 初始化
	scheduler := NewCronScheduler()

	_ = scheduler.AddTask(NewCronTask("@every 5s", testTask{id: 1, ratio: 0.5}.Execute))
	_ = scheduler.AddTask(NewCronTask("@every 2s", testTask{id: 2, ratio: 0.5}.Execute))
	_ = scheduler.AddTask(NewCronTask("@daily", testTask{id: 3, ratio: 0.5}.Execute))
	_ = scheduler.AddTask(NewCronTask("@0 */1 * * * ?s", testTask{id: 4, ratio: 0.5}.Execute))

	// 开始调度
	if err := scheduler.Execute(t.Context()); err != nil {
		t.Log(err)
		return
	}

	// 定时任务信息
	for _, cronTask := range scheduler.All() {
		fmt.Println(cronTask.GetMeta())
	}

	time.Sleep(10 * time.Second)

	for _, cronTask := range scheduler.All() {
		fmt.Println(cronTask.GetMeta())
	}
	select {}
}
