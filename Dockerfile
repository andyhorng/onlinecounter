FROM golang

ENV GOPATH /go/src
WORKDIR /go/src/app
COPY . /go/src/app

RUN go get -v -d
RUN go build

CMD ["./app", "8080"]
