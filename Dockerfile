FROM golang:1.12.3

COPY auction.go .
RUN go build auction.go
CMD ./auction