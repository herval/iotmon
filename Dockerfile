FROM golang:1.14

WORKDIR /src

ADD go.mod  /src
ADD go.sum /src
RUN go mod download

ADD . /src

RUN go build -o app ./bin

ENTRYPOINT ["/src/app"]


