FROM golang:1.20.1-alpine3.17 AS builder
WORKDIR /work

# Copy the Go Modules manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy the go source
COPY cmd ./cmd
COPY pkg ./pkg

RUN CGO_ENABLED=0 go build -o service cmd/main.go

FROM alpine:3.17.2
WORKDIR /bin
COPY --from=builder /work/service /bin/service
ENV GIN_MODE=debug
ENTRYPOINT ["/bin/service"]

EXPOSE 8080