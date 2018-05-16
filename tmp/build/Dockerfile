FROM alpine:3.6

RUN adduser -D pvc-operator
USER pvc-operator

ADD tmp/_output/bin/pvc-operator /usr/local/bin/pvc-operator
