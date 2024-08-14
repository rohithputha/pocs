package polling

import (
	"sync"
	"time"
)

type ec2 struct {
	status map[string]string
	mux    *sync.Mutex
}

func (e *ec2) getStatus(id string) string {
	e.mux.Lock()
	defer e.mux.Unlock()

	return e.status[id]
}

func (e *ec2) start(id string) {
	e.mux.Lock()
	if _, ok := e.status[id]; !ok {
		e.status[id] = "progress"
	} else {
		e.mux.Unlock()
		return
	}
	e.mux.Unlock()

	time.Sleep(10 * time.Second)

	e.mux.Lock()
	e.status[id] = "completed"
	e.mux.Unlock()
}
