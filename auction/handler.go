package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type BiddingResponse struct {
	AdId     string  `json:"adId"`
	BidPrice float64 `json:"bidPrice"`
	Err      error   `json:"-"`
}

var availableBidders = []string{"http://bidding:8080/buy/ad-slot", "http://bidding-1:8080/buy/ad-slot", "http://bidding-2:8080/buy/ad-slot"}

func HeaderBidHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the adPlacementId from path params
	adPlacementId := r.URL.Path[len("/bid/"):]
	log.Println("Request for Ad Placement", adPlacementId)

	var highestBid float64
	wg := sync.WaitGroup{}
	individualBids := make(chan BiddingResponse, len(availableBidders))
	defer close(individualBids)

	// initializing a context with 200ms timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(200*time.Millisecond))
	defer cancel()

	// Iterates over the available bidders and makes an http request to them in a goroutine
	for _, bidder := range availableBidders {
		wg.Add(1)
		go fetchAdBid(bidder, &wg, individualBids, ctx)
	}

	// extract the bidPrice from the channel and select the highest bid
	go func() {
		for bid := range individualBids {
			if highestBid < bid.BidPrice {
				highestBid = bid.BidPrice
			}
		}
	}()
	wg.Wait()

	// if no eligible bid is found, send a 204
	if highestBid == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	writeJSON(w, http.StatusOK, highestBid)
}

// Returns the bid response from the bid-service.
func fetchAdBid(bidder string, wg *sync.WaitGroup, bidCh chan BiddingResponse, ctx context.Context) {
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
		bidCh <- BiddingResponse{
			AdId:     response.AdId,
			BidPrice: response.BidPrice,
			Err:      err,
		}
	}
}
