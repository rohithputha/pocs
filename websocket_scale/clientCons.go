package websocket_scale

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"sync"
)

var companyIdSet map[string]interface{}
var chanIdConnMap map[string][]*websocket.Conn
var cidSetLock *sync.Mutex
var chanIdConnMapLock *sync.Mutex
var rds *redis.Client
var ctx = context.Background()

type redisMessage struct {
	ChanId string `json:"chanId"`
	Msg    string `json:"msg"`
}

func InitBackend() {
	companyIdSet = make(map[string]interface{})
	chanIdConnMap = make(map[string][]*websocket.Conn)
	cidSetLock = &sync.Mutex{}
	chanIdConnMapLock = &sync.Mutex{}
	rds = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
func addCompany(companyId string) {
	cidSetLock.Lock()
	defer cidSetLock.Unlock()

	if _, ok := companyIdSet[companyId]; ok {
		return
	}

	companyIdSet[companyId] = struct{}{}
	pubsub := rds.Subscribe(ctx, companyId)
	go rcvbrdcst(pubsub.Channel())
}

func rcvbrdcst(rdsMsg <-chan *redis.Message) {
	for msg := range rdsMsg {
		var redisMsg redisMessage
		json.Unmarshal([]byte(msg.Payload), &redisMsg)
		chanIdConnMapLock.Lock()
		if connList, ok := chanIdConnMap[redisMsg.ChanId]; ok {
			for _, conn := range connList {
				conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			}
		}
		chanIdConnMapLock.Unlock()
		// locking is very, very inefficient: should do fine for this POC though
	}
}

func addChanIdConn(chanId string, conn *websocket.Conn) {
	chanIdConnMapLock.Lock()
	defer chanIdConnMapLock.Unlock()

	if _, ok := chanIdConnMap[chanId]; !ok {
		chanIdConnMap[chanId] = make([]*websocket.Conn, 0)
	}
	chanIdConnMap[chanId] = append(chanIdConnMap[chanId], conn)
}
