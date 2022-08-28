package client

import (
	"github.com/gorilla/websocket"
)

type OccsClient struct {
	Id         string
	Connection *websocket.Conn
}
