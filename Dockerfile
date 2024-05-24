FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .

EXPOSE 8888

CMD ["./app"]
