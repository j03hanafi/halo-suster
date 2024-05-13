FROM golang:1.22.3 as builder
WORKDIR /app

ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o main cmd/api/main.go

FROM ubuntu:latest
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app/

COPY --from=builder /app/main .
COPY configs/config.toml ./configs/

EXPOSE 8080

ENV DB_HOST=host.docker.internal

CMD ./main