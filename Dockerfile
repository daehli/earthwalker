FROM golang:1.14.2-alpine

COPY . /opt/earthwalker/

WORKDIR /opt/earthwalker

RUN apk update && apk add --no-cache make && make

ENTRYPOINT ["./earthwalker"]
