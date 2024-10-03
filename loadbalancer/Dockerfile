FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o loadbalancer ./cmd/main.go

FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/loadbalancer .

CMD ["./loadbalancer"]