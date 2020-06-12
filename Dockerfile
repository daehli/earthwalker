FROM golang:1.14.2-alpine

COPY . /opt/earthwalker/

WORKDIR /opt/earthwalker

RUN go build

ENTRYPOINT ["./earthwalker"]
