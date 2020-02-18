FROM golang:onbuild

RUN mkdir /app
ADD . /app/
RUN go build -o auction .
CMD ["/app/auction"]