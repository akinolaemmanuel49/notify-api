package middlewares

import (
	"net/http"

	"github.com/akinolaemmanuel49/notify-api/utils"
)

func JWTAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := utils.ValidateJWT(w, r)
		if err != nil {
			utils.RespondWithError(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
