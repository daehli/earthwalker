FROM golang:1.14.2-alpine

COPY challenge /opt/earthwalker/challenge
COPY database /opt/earthwalker/database
COPY dynamicpages /opt/earthwalker/dynamicpages
COPY placefinder /opt/earthwalker/placefinder
COPY player /opt/earthwalker/player
COPY scores /opt/earthwalker/scores
COPY static /opt/earthwalker/static
COPY streetviewserver /opt/earthwalker/streetviewserver
COPY templates /opt/earthwalker/templates
COPY urlbuilder /opt/earthwalker/urlbuilder
COPY util /opt/earthwalker/util
COPY go.mod /opt/earthwalker/go.mod
COPY go.sum /opt/earthwalker/go.sum
COPY main.go /opt/earthwalker/main.go

WORKDIR /opt/earthwalker

RUN go build

ENTRYPOINT ["./earthwalker"]
