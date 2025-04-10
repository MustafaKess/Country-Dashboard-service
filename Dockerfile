FROM golang:1.22

LABEL maintainer="mail@host.tld"

COPY ./internal /go/src/app/cmd
COPY ./constants /go/src/app/constants
COPY ./go.mod /go/src/app/go.mod
COPY ./go.sum /go/src/app/go.sum
COPY ./main.go /go/src/app/main.go

WORKDIR /go/src/app

RUN CGO_ENABLED=0 GOOS=linux -a -ldflags '-extldflags "-static" ' -o server

EXPOSE 8080

CMD ["./server"]
