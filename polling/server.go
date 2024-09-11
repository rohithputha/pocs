package polling

import (
	"net/http"
	"sync"
)

func InitServer() {
	ec2 := &ec2{
		status: make(map[string]string),
		mux:    &sync.Mutex{},
	}
	poller := &poller{
		ec2: ec2,
	}

	http.HandleFunc("/shortpoll", poller.shortpoll)
	http.HandleFunc("/longpoll", poller.longpoll)
	http.ListenAndServe(":8080", nil)
}
