package middlewares

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

func Auth(w http.ResponseWriter, r *http.Request) error {
	log.Println("dumping headers...")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("%s: %s\n", name, value)
		}
	}

	secret := strings.TrimSpace(os.Getenv("API_TOKEN"))

	if r.Header.Get("Authorization") != "Bearer "+secret {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized"))
		return errors.New("unauthorized")
	}

	return nil
}
