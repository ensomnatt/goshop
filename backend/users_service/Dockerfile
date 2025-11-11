FROM golang:1.25.3 AS builder

WORKDIR /app

COPY go.sum go.mod ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/users_service/main.go

FROM debian:bookworm-slim

WORKDIR /root

COPY --from=builder /app/app .

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

CMD ["./app"]
