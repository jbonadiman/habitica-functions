package middlewares

import (
	"net/http"
	"os"
	"strings"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	secret := strings.TrimSpace(os.Getenv("API_TOKEN"))

	if r.Header.Get("Authorization") != "Bearer "+secret {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized"))
		return
	}
}
