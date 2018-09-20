FROM golang:1.9-alpine

ADD . /go/src/github.com/banzaicloud/pvc-operator
WORKDIR /go/src/github.com/banzaicloud/pvc-operator
RUN go build -o /tmp/pvc-operator cmd/pvc-operator/main.go

FROM alpine:3.6

COPY --from=0 /tmp/pvc-operator /usr/local/bin/pvc-operator
RUN apk update && apk add ca-certificates
RUN adduser -D pvc-operator

USER pvc-operator

ENTRYPOINT ["/usr/local/bin/pvc-operator"]
