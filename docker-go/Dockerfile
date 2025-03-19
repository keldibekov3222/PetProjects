FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0


RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]
