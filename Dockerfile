FROM golang:alpine

RUN mkdir /src
COPY ./ /src
WORKDIR /src

RUN go build -o /src/parser ./cmd/parser/main.go
ENTRYPOINT ["/src/parser"]