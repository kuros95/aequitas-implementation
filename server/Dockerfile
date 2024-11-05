ARG PORT=2220

FROM golang:1.23 AS builder

WORKDIR /builder

RUN mkdir server && mkdir stayalive

COPY ./server/main.go /builder/server

COPY /stayalive/ /builder/stayalive

COPY go.mod go.sum /builder/

RUN go build ./server/main.go

FROM alpine:3.20 AS runner

WORKDIR /runner

COPY --from=builder /builder/main runner/

EXPOSE 2220 2222 2224

CMD [ "./main", "-p", "PORT" ]