FROM golang:1.23.4-alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/dnsupd/
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o /go/bin/dnsupd

RUN go install github.com/bugsnag/panic-monitor@latest

FROM scratch
ARG GIT_COMMIT
ENV GIT_COMMIT ${GIT_COMMIT}
ENV BUGSNAG_SOURCE_ROOT=/go/src/dnsupd/

COPY --from=builder /go/bin/dnsupd /go/bin/dnsupd
COPY --from=builder /go/bin/panic-monitor /go/bin/panic-monitor

EXPOSE 80
ENTRYPOINT ["/go/bin/dnsupd"]
