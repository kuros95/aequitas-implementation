FROM golang:1.23 AS builder

WORKDIR /builder

RUN mkdir client && mkdir stayalive && mkdir aequitas 

RUN apt update && apt install iproute2 -y 

COPY ./client/main.go /builder/client

COPY /stayalive/ /builder/stayalive

COPY /utils/ /builder/utils

COPY go.mod go.sum tc-on-host.sh /builder/

RUN go build -o client ./client/

EXPOSE 2220 2222 2224

ENTRYPOINT [ "./client/client" ]