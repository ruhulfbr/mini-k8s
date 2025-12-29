package main

import (
	"log"

	"github.com/ruhulfbr/mini-k8s/cmd/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal("Application run error", "err", err)
	}
}
