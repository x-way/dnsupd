FROM golang:1.16.1-alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/dnsupd/
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o /go/bin/dnsupd

FROM scratch
COPY --from=builder /go/bin/dnsupd /go/bin/dnsupd

EXPOSE 80
ENTRYPOINT ["/go/bin/dnsupd"]
