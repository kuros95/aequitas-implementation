FROM golang:1.23 AS builder

WORKDIR /builder

RUN mkdir client && mkdir stayalive && mkdir aequitas

COPY ./client/main.go /builder/client

COPY /stayalive/ /builder/stayalive

COPY /aequitas/ /builder/aequitas

COPY go.mod go.sum /builder/

RUN go build ./client/main.go

FROM alpine:3.20 AS runner

RUN apk add iproute2

WORKDIR /runner

COPY --from=builder /builder/main runner/

COPY tc-on-host.sh /runner/

RUN src_ip=$(ip addr show eth0 scope global | grep 172 | awk '{print $2}' | cut -d "/" -f 1) && ./tc-on-host.sh

EXPOSE 2220 2222 2224

CMD [ "./main" ]