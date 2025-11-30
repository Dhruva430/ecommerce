package main

import (
	"api/configs"
	"api/internals/routes"
	"api/models/db"
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func loadEnvVariables() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("⚠️ error loading .env file: %w", err)
	}
	return nil
}

func connectDB() *sql.DB {
	dbURL := configs.GetDBURI()

	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}
	log.Println("✅ Database connected")
	return sqlDB
}

func main() {
	loadEnvVariables()

	conn := connectDB()
	queries := db.New(conn)
	g := routes.SetupRouter(queries, conn)
	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
