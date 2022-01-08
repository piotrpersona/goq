FROM golang:1.18beta1-alpine

COPY go.* .

RUN go mod download

COPY . .