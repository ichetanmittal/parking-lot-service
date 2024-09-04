FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /parking-lot-service ./cmd/server

EXPOSE 8080

CMD ["/parking-lot-service"]