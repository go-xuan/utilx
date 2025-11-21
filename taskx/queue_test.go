package taskx

import (
	"context"
	"fmt"
	"testing"
)

func TestQueueScheduler(t *testing.T) {
	queue := NewQueue("queue_test")
	queue.Add("task1", testTask{1, 0.95}.Execute)                // 正常插入 task1
	queue.Add("task2", testTask{2, 0.95}.Execute)                // 默认尾插 task1->task2
	queue.AddAfter("task1", "task3", testTask{3, 0.95}.Execute)  // 插队到task1后面 task1->task3->task2
	queue.Add("task4", testTask{4, 0.95}.Execute)                // 默认尾插 task1->task3->task2->task4
	queue.AddBefore("task4", "task5", testTask{5, 0.95}.Execute) // 插队到task4前面 task1->task3->task2->task5->task4
	queue.AddTail("task6", testTask{6, 0.95}.Execute)            // 插队到末位 task1->task3->task2->task5->task4->task6
	queue.AddHead("task7", testTask{7, 0.95}.Execute)            // 插队到首位 task7->task1->task3->task2->task5->task4->task6
	queue.Remove("task2")                                        // 删除task2 task7->task1->task3->task5->task4->task6
	fmt.Println(queue.GetNames())
	if err := queue.Execute(context.TODO()); err != nil {
		t.Log(err)
	}
}
