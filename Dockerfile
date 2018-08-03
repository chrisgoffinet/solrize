FROM golang:1.11beta2 as builder
ENV WORKDIR /go/src/github.com/chrisgoffinet/solrize
WORKDIR ${WORKDIR}

COPY . ${WORKDIR}
ENV GO111MODULE on
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app

FROM alpine:latest
ENV WORKDIR /go/src/github.com/chrisgoffinet/solrize
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder ${WORKDIR}/app /go/bin/app
ENTRYPOINT [ "/go/bin/app" ]