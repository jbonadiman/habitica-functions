package habitica_functions

import (
	"fmt"
	"net/http"

	"habitica_functions/internal"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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

	err := internal.FinishFocusSession()
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

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Cache-Control", "s-maxage=3, stale-while-revalidate=59")
	w.WriteHeader(http.StatusOK)
}
