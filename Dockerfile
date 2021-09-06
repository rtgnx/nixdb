FROM golang:1.16.4-alpine3.13

RUN apk update && apk add alpine-sdk linux-pam linux-pam-dev

RUN mkdir -p /go/src/github.com/Reverse-Labs/nixdb
WORKDIR /go/src/github.com/Reverse-Labs/nixdb
COPY . .
RUN go fmt
RUN addgroup testuser
RUN adduser -G testuser -D testuser
RUN passwd -d testuser
RUN go test

RUN go build -o /usr/bin/nixdbd cmd/nixdbd/*.go
RUN /usr/bin/nixdbd genkey --secretKey /etc/nixdb.key --secretSize 4096

CMD ["/usr/bin/nixdbd", "serve", "--baseDir", "/etc", "--minGID", "0", "--minUID", "0"]