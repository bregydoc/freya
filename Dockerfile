FROM golang:1.11
WORKDIR /go/src/freya
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["freya"]