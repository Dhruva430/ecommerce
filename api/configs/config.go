package configs

import "os"

func GetDBURI() string {
	return os.Getenv("DATABASE_URL")
}
