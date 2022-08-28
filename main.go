package main

import (
	"github.com/Erik142/occs-server/internal/server"
	"log"
)

func main() {
	log.Default().Println("Hello World")
	server.Run(8888)
}
