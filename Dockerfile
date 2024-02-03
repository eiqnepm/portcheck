FROM golang:1.20.2-alpine3.17 AS build

WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd

RUN go mod download

RUN go build -ldflags="-s -w" -o /usr/local/bin/app cmd/portcheck/main.go

FROM alpine:3.17

RUN apk update
# RUN apk add --no-cache docker

COPY --from=build /usr/local/bin/app /app

CMD ["/app"]
