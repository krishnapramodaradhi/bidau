FROM golang:1.22.1-alpine3.19 AS build-stage

WORKDIR /app

COPY go.mod ./

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /bidding-service ./main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /bidding-service /bidding-service

USER nonroot:nonroot

ENTRYPOINT ["/bidding-service"]
