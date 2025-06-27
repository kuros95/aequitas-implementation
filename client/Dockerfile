FROM golang:1.23 AS builder

WORKDIR /builder

RUN mkdir client && mkdir sendmess && mkdir aequitas 

RUN apt-get update && apt-get install iproute2 tcpdump -y 

COPY ./client/main.go /builder/client

COPY /sendmess/ /builder/sendmess

COPY /utils/ /builder/utils

COPY go.mod go.sum tc-on-host.sh run-tcpdump.sh 32KB-payload 64KB-payload /builder/

WORKDIR /builder/client

RUN go build

EXPOSE 2220 2222 2224

WORKDIR /builder

ENTRYPOINT [ "./client/client" ]