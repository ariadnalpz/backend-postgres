package utils

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Fatal("El archivo .env no existe")
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando archivo .env")
	}
}