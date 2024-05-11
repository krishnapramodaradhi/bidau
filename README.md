# bidau (short for Bidders & Auctioneers)
A Go based microservices application where bidders can place their bid for an Ad Space using the bidding-service. An Auctioneer can sell the Ad space to the highest bidder using the aution-service.

## Run it locally
### Prerequisites: [Docker Desktop](https://www.docker.com/products/docker-desktop/) should be installed.
1. Clone the repository.
2. Navigate to the cloned directory.
3. Run ```make up``` to start the services.
4. Access the auction-service at ```http://localhost:8084/bid/<adPlacementId>``` to get the highest bid.
