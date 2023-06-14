package main

import (
	"log"

	"github.com/capcom6/swarm-gateway-tutorial/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
