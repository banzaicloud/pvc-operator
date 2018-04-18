FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -y ca-certificates

ADD tmp/_output/bin/pvc-handler /usr/local/bin/pvc-handler

RUN useradd -ms /bin/bash pvc-handler
USER pvc-handler
