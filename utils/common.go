package utils

import (
	"net/http"

	"github.com/akinolaemmanuel49/notify-api/config"
)

var cfg config.Config
var privateKey = []byte(cfg.JWT.KEY)

func RespondWithError(w http.ResponseWriter, errorMessage string, code int) {
	http.Error(w, errorMessage, code)
}
