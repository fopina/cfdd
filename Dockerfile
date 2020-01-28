FROM alpine:3.11 as certs

RUN apk add --no-cache ca-certificates

FROM golang:1.13-alpine as builder

RUN apk add --no-cache make

WORKDIR /go/src/app

ADD go.mod /go/src/app
ADD go.sum /go/src/app
RUN go mod download

ADD . /go/src/app

ARG VERSION=dev
RUN make build

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/app/dist/cfdd /cfdd

ARG VERSION=dev
LABEL version="${VERSION}" maintainer="fopina <https://github.com/fopina/cfdd/>"

ENTRYPOINT [ "/cfdd" ]
