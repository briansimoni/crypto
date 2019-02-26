FROM golang:1.10

ADD . /app/go/src/github.com/briansimoni/crypto

WORKDIR /app/go/src/github.com/briansimoni/crypto

RUN go build -o crypto

CMD /app/go/src/github.com/briansimoni/crypto/crypto