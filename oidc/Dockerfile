FROM golang:1.24 AS builder
WORKDIR /build
COPY go.mod go.sum . 
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server/main.go

FROM alpine:latest
ENV CRON="*/5 * * * *"
WORKDIR /app
COPY --from=builder /build/main /app/main
EXPOSE 8080
CMD ["/app/main"]
