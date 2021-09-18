package task

import (
	"github.com/jixindatech/pkg/config"
	"github.com/jixindatech/pkg/core/alert"
	"github.com/jixindatech/pkg/queue"
	"github.com/jixindatech/pkg/webserver/models"
	"time"
)

type TaskMysql struct {
	Name  string
	Queue queue.Queue
	Alert alert.Alert
}

func (t *TaskMysql) GetName() string {
	return t.Name
}

func (t *TaskMysql) Run() {
	ch, _ := t.Queue.GetHandler()
	sqlChannel := ch.(chan config.SqlMsg)

	for true {
		select {
		case res := <-sqlChannel:
			//write to es
			go models.SaveEvent(res)
			if res.Alert > 0 {
				msg := res.Name + " 触发告警：" + res.Sql
				go t.Alert.Send("数据库审计", msg, res.Name)
			}
		case <-time.After(time.Second * 3):
			//: TODO
		}

	}
}
