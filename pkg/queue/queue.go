package queue

type Queue interface {
	Init(int) error
	Put(interface{}) error
	GetHandler() (interface{}, error)
}

func GetQueue(name string, num int) (Queue, error) {
	if name == "chan" {
		queueInstance := &QueueChan{}
		err := queueInstance.Init(num)
		if err != nil {
			return nil, err
		}
		return queueInstance, nil
	}
	return nil, nil
}
