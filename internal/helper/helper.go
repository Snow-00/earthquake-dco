package helper

import (
	"log"
	"net/http"
)

func ResultCheck(w http.ResponseWriter, text string, status ...int) {
	statusCode := http.StatusNoContent

	if len(status) > 0 {
		statusCode = status[0]
	}

	log.Println(text)
	w.WriteHeader(statusCode)
}
