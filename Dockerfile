FROM golang:1.13

WORKDIR /usr/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["main"]
