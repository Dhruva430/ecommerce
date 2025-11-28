package configs

import (
	"log"

	"github.com/Dhruva430/ecommerce/prisma/db"
)

var Prisma *db.PrismaClient

func InitPrisma() {
	Prisma = db.NewClient()
	if err := Prisma.Prisma.Connect(); err != nil {
		log.Fatal("❌ Failed to connect to DB:", err)
	}
	log.Println("✅ Prisma connected")
}
