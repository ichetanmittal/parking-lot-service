FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /parking-lot-service ./cmd/server

EXPOSE 8080

CMD [ "/parking-lot-service" ]