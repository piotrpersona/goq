FROM golang:1.18beta1-bullseye

WORKDIR /goq

COPY go.* .

RUN go mod download

COPY . .