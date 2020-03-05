FROM golang:1.13

RUN mkdir -p $GOPATH/src/app
WORKDIR /$GOPATH/src/app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["main"]
