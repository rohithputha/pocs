package websocket_scale

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type clientSetup struct {
	CompanyId string `json:"companyId"`
	ChannelId string `json:"channelId"` //should actually be a slice of strings. should be fine for this POC.
}

func wsClient(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	//defer conn.Close()
	fmt.Println("Client connected")
	_, message, err := conn.ReadMessage()
	if err != nil {
		return
	}
	var cs clientSetup
	err = json.Unmarshal(message, &cs)
	fmt.Println(cs)
	fmt.Println(err)
	if err != nil {
		return
	}
	addCompany(cs.CompanyId)
	addChanIdConn(cs.ChannelId, conn)
}

func InitServer() {
	router := http.NewServeMux()
	router.HandleFunc("/ws", wsClient)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
}
