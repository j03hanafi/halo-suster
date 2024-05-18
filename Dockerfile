FROM golang:1.22.3 as builder
WORKDIR /app

ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v3 go build -ldflags='-s -w' -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add dumb-init
# RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app/

COPY --from=builder /app/main .
COPY configs/config.toml.example ./configs/config.toml

EXPOSE 8080

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ./main