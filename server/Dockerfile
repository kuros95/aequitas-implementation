FROM golang:1.23 AS builder

WORKDIR /builder

RUN mkdir server && mkdir stayalive

COPY ./server/main.go /builder/server

COPY /stayalive/ /builder/stayalive

COPY go.mod go.sum /builder/

RUN pwd && ls -al && go build ./server/main.go

FROM alpine:3.20 AS runner

WORKDIR /server

COPY --from=builder /builder/main runner/

CMD [ "./main" ]