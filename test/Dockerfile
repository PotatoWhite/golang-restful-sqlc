FROM golang:1.20.1-alpine3.17
WORKDIR /root

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg

CMD CGO_ENABLED=0 go test -cover ./...