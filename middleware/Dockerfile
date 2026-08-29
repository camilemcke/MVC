FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o middleware .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/middleware .
EXPOSE 8082
CMD ["./middleware"]
