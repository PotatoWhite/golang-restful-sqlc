FROM golang:1.20.1-alpine3.17 AS builder
WORKDIR /work

# Copy the Go Modules manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy the go source
COPY main.go ./