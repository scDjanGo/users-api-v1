FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6

COPY . .
RUN swag init -g main.go -o src/docs
RUN CGO_ENABLED=0 go build -o server .

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]