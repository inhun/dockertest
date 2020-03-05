FROM golang:1.13

WORKDIR /go/src/app
COPY . .

ENV PATH $GOPATH/bin:/opt/go/bin:$PATH
ENV GOPATH /go

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["main"]
