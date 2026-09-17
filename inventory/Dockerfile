FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o inventory .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/inventory .

EXPOSE 8080

CMD ["./inventory"]