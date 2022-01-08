FROM golang:1.18beta1-bullseye

WORKDIR /goq

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .