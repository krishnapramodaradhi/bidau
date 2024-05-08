package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	random "math/rand"
	"net/http"
	"strconv"
)

const PORT = ":8443"

type BiddingResponse struct {
	AdId     string  `json:"adId"`
	BidPrice float64 `json:"bidPrice"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/buy/ad-slot", func(w http.ResponseWriter, r *http.Request) {
		if rInt := random.Intn(10); rInt != 0 && rInt%4 == 0 {
			log.Println(rInt)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		br := BiddingResponse{AdId: generateRandomId(), BidPrice: truncate(random.Float64() * 100)}
		writeJSON(w, http.StatusOK, br)
	})

	log.Println("Listening on port", PORT)
	log.Fatal(http.ListenAndServe(PORT, mux))
}

func truncate(num float64) float64 {
	truncatedString := fmt.Sprintf("%.2f", num)
	truncatedNum, _ := strconv.ParseFloat(truncatedString, 64)
	return float64(truncatedNum)
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
