FROM golang:1-alpine AS build

WORKDIR /usr/src/portcheck

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN go build -ldflags="-s -w" -o /usr/local/bin/portcheck cmd/portcheck/main.go

FROM alpine:3.17

COPY --from=build /usr/local/bin/portcheck /portcheck

CMD ["/portcheck"]
