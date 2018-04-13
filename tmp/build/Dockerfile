FROM alpine:3.6

ADD tmp/_output/bin/pvc-handler /usr/local/bin/pvc-handler

RUN adduser -D pvc-handler
USER pvc-handler
