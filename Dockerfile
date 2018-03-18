# Build Geth in a stock Go builder container
FROM golang:1.10-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

ADD . /go-ubiq
RUN cd /go-ubiq && make gubiq

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-ubiq/build/bin/gubiq /usr/local/bin/

EXPOSE 8588 8589 30303 30303/udp 30304/udp
ENTRYPOINT ["gubiq"]
