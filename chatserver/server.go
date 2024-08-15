package chatserver

import (
	"encoding/json"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"sync"
)

var connMap map[string]socketio.Conn
var mapLock *sync.Mutex

type lockType struct {
	connMapLock bool
}

type message struct {
	From string `json:"from"`
	To   string `json:"to"`
	Msg  string `json:"msg"`
}

type lockOption func(*lockType)

func connMapLock() lockOption {
	return func(lt *lockType) {
		lt.connMapLock = true
	}
}

func withLock(f func(), lockOpts ...lockOption) {
	lt := &lockType{}
	for _, opt := range lockOpts {
		opt(lt)
	}
	if lt.connMapLock {
		mapLock.Lock()
		defer mapLock.Unlock()
	}
	f()
}

func InitServer() {
	server := socketio.NewServer(nil)
	connMap = make(map[string]socketio.Conn)
	mapLock = &sync.Mutex{}
	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "register", func(s socketio.Conn, userId string) {
		withLock(func() {
			fmt.Println("register", userId, s)
			connMap[userId] = s
			fmt.Println(connMap)
		}, connMapLock())
	})

	server.OnEvent("/", "msg", func(s socketio.Conn, msg string) {
		withLock(func() {
			m := &message{}
			fmt.Println("msg", msg)
			err := json.Unmarshal([]byte(msg), &m)
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Println("msg", m)
			if conn, ok := connMap[m.To]; ok {
				fmt.Println("found connection " + m.To)
				fmt.Println(connMap)
				conn.Emit("msg", m.Msg)
			}
		}, connMapLock())
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		withLock(func() {
			for k, v := range connMap {
				if v == s {
					delete(connMap, k)
				}
			}
		}, connMapLock())
	})
	go server.Serve()
	defer server.Close()
	http.Handle("/", server)
	//http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
