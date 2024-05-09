package main

import (
	"io"
	"log"
	"net/http"
)

const bidPriceLowerThreshold float64 = 29.99

var availableBidders = []string{"http://bidding:8080/buy/ad-slot", "http://bidding:8081/buy/ad-slot", "http://bidding:8082/buy/ad-slot", "http://bidding:8083/buy/ad-slot"}

func HeaderBidHandler(w http.ResponseWriter, r *http.Request) error {
	adPlacementId := r.URL.Path[len("/bid/"):]
	log.Println("Request for Ad Placement", adPlacementId)
	// var highestBid float64
	for _, v := range availableBidders {
		res, err := http.Get(v)
		if err != nil {
			log.Println("An error occured in handler", err)
		}
		if res.StatusCode == http.StatusOK {
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				log.Println("an error occured while reading the response body", err)
			}
			log.Println(resBody)

		}
	}
	return writeJSON(w, http.StatusOK, map[string]string{"message": "Done"})
}
