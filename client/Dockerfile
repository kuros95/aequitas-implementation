FROM golang:1.23 AS builder

WORKDIR /builder

RUN mkdir client && mkdir stayalive && mkdir aequitas

COPY ./client/main.go /builder/client

COPY /stayalive/ /builder/stayalive

COPY /aequitas/ /builder/aequitas

COPY tc-on-host.sh /runner/

COPY go.mod go.sum /builder/

RUN go build -o client ./client/

# FROM alpine:3.20 AS runner

# RUN apk add iproute2

# WORKDIR /runner
 
# COPY --from=builder /builder/client/client /runner/



EXPOSE 2220 2222 2224

ENTRYPOINT [ "./tc-on-host.sh" ]

CMD [ "/bin/sh", "./client" ]