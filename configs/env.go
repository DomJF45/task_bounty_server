package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env")
	}

	return os.Getenv("MONGO_URI")
}

func EnvJWTSecret() []byte {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env")
	}

	return []byte(os.Getenv("JWT_SECRET"))
}
