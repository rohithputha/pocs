package polling

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type poller struct {
	ec2 *ec2
}

func (p *poller) shortpoll(w http.ResponseWriter, r *http.Request) {
	go p.ec2.start(r.URL.Query().Get("id"))
	fmt.Println(p.ec2.getStatus(r.URL.Query().Get("id")))
	io.WriteString(w, p.ec2.getStatus(r.URL.Query().Get("id")))
}

func (p *poller) longpoll(w http.ResponseWriter, r *http.Request) {
	go p.ec2.start(r.URL.Query().Get("id"))

	for {
		time.Sleep(1 * time.Second)
		fmt.Println(p.ec2.getStatus(r.URL.Query().Get("id")))
		if status := p.ec2.getStatus(r.URL.Query().Get("id")); status == "completed" {
			io.WriteString(w, status)
			return
		}
	}
	return
}
