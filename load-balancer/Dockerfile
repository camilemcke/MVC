FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY . .

RUN go build -o load-balancer .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/load-balancer .

EXPOSE 8080

CMD ["./load-balancer"]