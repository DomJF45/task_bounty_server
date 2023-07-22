package configs

import (
	"log"
	"os"
)

func EnvMongoURI() string {
	/*
		err := godotenv.Load()
		if err != nil {
			log.Fatal("error loading env")
		}
	*/

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set in environment")
	}

	return mongoURI
}

func EnvJWTSecret() []byte {
	/*
		err := godotenv.Load()
		if err != nil {
			log.Fatal("error loading env")
		}
	*/
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set in environment")
	}

	return []byte(jwtSecret)
}
