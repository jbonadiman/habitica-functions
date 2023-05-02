package habitica_functions

import (
	"fmt"
	"net/http"

	"habitica_functions/internal"
	"habitica_functions/internal/middlewares"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	err := middlewares.Auth(w, r)
	if err != nil {
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write(
			[]byte(fmt.Sprintf(
				"the method %q is not allowed",
				r.Method,
			)),
		)
		return
	}

	err = internal.FinishFocusSession()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(
			[]byte(fmt.Sprintf(
				"error on habitica integration: %v",
				err,
			)),
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}
