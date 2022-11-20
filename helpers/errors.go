package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func CheckErr(err error) {
	if err != nil {
		log.Printf("message: %v", err)
	}
}

func ErrorResonse(msg string, status int, res http.ResponseWriter) {
	if msg == "" {
		return
	}
	log.Printf("error: %v", msg)

	res.Header().Add("content-type", "application/json")
	res.WriteHeader(status)
	payload, err := json.Marshal(map[string]string{"error": msg})
	CheckErr(err)
	res.Write(payload)
}

func SeverError(err error, res http.ResponseWriter) {
	if err == nil {
		return
	}
	log.Printf("error: %v", err)

	res.Header().Add("content-type", "application/json")
	res.WriteHeader(http.StatusInternalServerError)
	payload, err2 := json.Marshal(map[string]string{"error": err.Error()})
	CheckErr(err2)
	res.Write(payload)
}
