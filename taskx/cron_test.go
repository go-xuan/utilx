package taskx

import (
	"fmt"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	// 初始化
	scheduler := NewCronScheduler()

	_ = scheduler.Add(NewCronTask("task1", "@every 5s", &testTask{id: 1, ratio: 0.6}))
	_ = scheduler.Add(NewCronTask("task2", "@every 2s", &testTask{id: 2, ratio: 0.5}))
	_ = scheduler.Add(NewCronTask("task3", "@daily", &testTask{id: 3, ratio: 0.6}))
	_ = scheduler.Add(NewCronTask("task4", "@0 */1 * * * ?s", &testTask{id: 4, ratio: 0.6}))

	// 开始调度
	if err := scheduler.Execute(t.Context()); err != nil {
		t.Fatal(err)
		return
	}

	// 定时任务信息
	for _, task := range scheduler.All() {
		fmt.Println(task.GetMeta())
	}

	time.Sleep(20 * time.Second)

	for _, task := range scheduler.All() {
		fmt.Println(task.GetMeta())
	}
	select {}
}
