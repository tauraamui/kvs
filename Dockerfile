FROM golang:1.20

WORKDIR /go/src/app

COPY . .

WORKDIR /go/src/app/v2

RUN go mod download
RUN go test -c ./query

ENTRYPOINT ["./query.test", "-test.bench=."]
