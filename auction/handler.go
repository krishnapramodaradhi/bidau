package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

const bidPriceLowerThreshold float64 = 29.99

type BiddingResponse struct {
	AdId     string  `json:"adId"`
	BidPrice float64 `json:"bidPrice"`
}

var availableBidders = []string{"http://bidding:8080/buy/ad-slot", "http://bidding-1:8080/buy/ad-slot", "http://bidding-2:8080/buy/ad-slot"}

func HeaderBidHandler(w http.ResponseWriter, r *http.Request) {
	adPlacementId := r.URL.Path[len("/bid/"):]
	log.Println("Request for Ad Placement", adPlacementId)

	var highestBid float64
	wg := sync.WaitGroup{}
	individualBids := make(chan float64, len(availableBidders))
	defer close(individualBids)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(200*time.Millisecond))
	defer cancel()

	for _, bidder := range availableBidders {
		wg.Add(1)
		go fetchAdBid(bidder, &wg, individualBids, ctx)
	}
	go func() {
		for bidPrice := range individualBids {
			if highestBid < bidPrice {
				highestBid = bidPrice
			}
		}
	}()
	wg.Wait()
	if highestBid == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	writeJSON(w, http.StatusOK, highestBid)
}

func fetchAdBid(bidder string, wg *sync.WaitGroup, bidCh chan float64, ctx context.Context) {
	var response BiddingResponse
	defer wg.Done()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bidder, nil)
	if err != nil {
		log.Println("An error occured in handler", err)
	}
	res, err := http.DefaultClient.Do(req)
	select {
	case <-ctx.Done():
		log.Println("operation timedout as deadline exceeded")
		return
	default:
		log.Println("Operation in progress")
	}
	if res.StatusCode == http.StatusOK {
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			log.Println("An error occured in handler", err)
		}
		bidCh <- response.BidPrice
	}
}
