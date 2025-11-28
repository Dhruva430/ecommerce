package main

import (
	"log"

	"github.com/Dhruva430/ecommerce/prisma/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	_ = godotenv.Load()

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return err
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			log.Println("Error disconnecting:", err)
		}
	}()

	gin.SetMode(gin.DebugMode) // change to gin.ReleaseMode in production

	r := gin.Default()

	r.SetTrustedProxies([]string{"127.0.0.1"}) // recommended for dev

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ecommerce backend running"})
	})

	log.Println("Server running on http://localhost:8080")
	return r.Run(":8080")
}
