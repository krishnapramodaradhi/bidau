package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
)

const PORT = ":8443"

type HttpFuncHandler func(w http.ResponseWriter, r *http.Request) error

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/bid/{adPlacementId}", makeHttp(HeaderBidHandler))

	log.Println("Listening on port", PORT)
	log.Fatal(http.ListenAndServe(PORT, mux))
}

func makeHttp(f HttpFuncHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Println(err)
		}
	}
}

func generateRandomId() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Println("Error while generating a random id: ", err)
		return ""
	}
	return hex.EncodeToString(b)
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}
