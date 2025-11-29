package main

import (
	"api/internals/routes"
	"log"
)

func main() {
	g := routes.SetupRouter()
	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
