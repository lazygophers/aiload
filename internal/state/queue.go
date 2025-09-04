package state

import (
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/aiload/internal/memqueue"
)

var (
	MqUpdateModel       *memqueue.Queue[*aiload.MqTaskId]
	MqCheckChannelModel *memqueue.Queue[*aiload.MqTaskId]
)

func LoadMemQueue() {
	MqUpdateModel = memqueue.New[*aiload.MqTaskId](&memqueue.Config[*aiload.MqTaskId]{
		ConcurrentCount: 1,
		MaxRetry:        1,
	})

	MqCheckChannelModel = memqueue.New[*aiload.MqTaskId](&memqueue.Config[*aiload.MqTaskId]{
		ConcurrentCount: 1,
		MaxRetry:        1,
	})
}
