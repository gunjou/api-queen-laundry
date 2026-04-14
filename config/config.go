package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
	Port  string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	return Config{
		DBUrl: os.Getenv("DATABASE_URL"),
		Port:  os.Getenv("PORT"),
	}
}