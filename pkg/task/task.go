package task

import (
	"github.com/jixindatech/pkg/core/alert"
	"github.com/jixindatech/pkg/queue"
)

type Task interface {
	GetName() string
	Run()
}

func GetTask(taskType string, name string, queueInstance queue.Queue, alertInstance alert.Alert) (Task, error) {
	if taskType == "mysql" {
		task := &TaskMysql{
			Name:  name,
			Queue: queueInstance,
			Alert: alertInstance,
		}

		return task, nil
	}
	return nil, nil
}
