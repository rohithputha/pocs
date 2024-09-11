package main

import (
	"pocs/websocket_scale"
)

//
//func main() {
//	polling.InitServer()
//}

//func main() {
//	chatserver.InitServer()
//}

func main() {
	websocket_scale.InitBackend()
	websocket_scale.InitServer()
}
