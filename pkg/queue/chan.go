package queue

import (
	"github.com/jixindatech/sqlaudit/pkg/config"
)

type QueueChan struct {
	num   int
	Queue chan config.SqlMsg
}

func (c *QueueChan) Init(num int) error {
	c.num = num
	c.Queue = make(chan config.SqlMsg, num)
	return nil
}

func (c *QueueChan) Put(msg interface{}) error {
	c.Queue <- msg.(config.SqlMsg)
	return nil
}

func (c *QueueChan) GetHandler() (interface{}, error) {
	//save := f.(callback)
	return c.Queue, nil
}
