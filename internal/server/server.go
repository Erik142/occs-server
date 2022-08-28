package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Erik142/occs-server/internal/client"
	"github.com/Erik142/occs-server/internal/message"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type OccsServer struct {
	clients []*client.OccsClient
}

var server = &OccsServer{}

func (s *OccsServer) ProcessMessage(client *client.OccsClient, messageType int, data []byte) {
	m := message.Message{}

	if err := json.Unmarshal(data, &m); err != nil {
		log.Default().Println("An error occured!")
		return
	}

	switch m.Type {
	case message.Subscribe:
		log.Default().Println(fmt.Sprintf("Subscribe: %s", m.Data))
		s.Send(client, fmt.Sprintf("{\"response\": \"%s\"}", m.Data))
		break
	case message.Unsubscribe:
		log.Default().Println(fmt.Sprintf("Unsubscribe: %s", m.Data))
		s.Send(client, fmt.Sprintf("{\"response\": \"%s\"}", m.Data))
		break
	case message.Publish:
		log.Default().Println(fmt.Sprintf("Publish: %s", m.Data))
		s.Send(client, fmt.Sprintf("{\"response\": \"%s\"}", m.Data))
		break
	default:
		log.Default().Println(fmt.Sprintf("Unrecognized message type: %s", m.Data))
		s.Send(client, fmt.Sprintf("{\"response\": \"%s\"}", m.Data))
		break
	}
}

func (s *OccsServer) AddClient(client *client.OccsClient) {
	for _, c := range s.clients {
		if c.Id == client.Id {
			return
		}
	}

	s.clients = append(s.clients, client)
}

func (s *OccsServer) RemoveClient(client *client.OccsClient) {
	for i, c := range s.clients {
		if c.Id == client.Id {
			copy(s.clients[i:], s.clients[i+1:])
			s.clients[len(s.clients)-1] = nil // or the zero value of T
			s.clients = s.clients[:len(s.clients)-1]
			break
		}
	}
}

func (s *OccsServer) Send(client *client.OccsClient, message string) {
	client.Connection.WriteMessage(1, []byte(message))
}

func webSocketHandler(ctx *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// upgrades connection to websocket
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// create new client & add to client list
	client := client.OccsClient{
		Id:         uuid.Must(uuid.NewRandom()).String(),
		Connection: conn,
	}

	// greet the new client
	server.Send(&client, "Server: Welcome! Your ID is "+client.Id)

	// message handling
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			server.RemoveClient(&client)
			return
		}
		server.ProcessMessage(&client, messageType, p)
	}
}

func Run(port int) {
	log.Default().Println("Creating new server")
	router := gin.Default()
	router.GET("/ws", webSocketHandler)
	router.Run(fmt.Sprintf(":%d", port))
}
