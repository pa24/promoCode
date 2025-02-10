FROM golang:1.23.2-alpine

WORKDIR /app

COPY . .

COPY go.mod go.sum ./
RUN go mod tidy

RUN go build -o promoCode ./cmd/server/main.go

CMD ["./promoCode"]