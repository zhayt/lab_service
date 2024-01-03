FROM golang:latest AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

CMD ["./app"]
