package configs

import "os"

func GetDBURI() string {
	return os.Getenv("DATABASE_URL")
}

func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}
