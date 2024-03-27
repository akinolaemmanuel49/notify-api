package utils

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/notify-api/config"
	"github.com/joho/godotenv"
)

var cfg config.Config
var privateKey = []byte(cfg.JWT.Key)

func RespondWithError(w http.ResponseWriter, errorMessage string, code int) {
	http.Error(w, errorMessage, code)
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}
