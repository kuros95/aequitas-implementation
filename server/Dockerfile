ARG PORT=2220

FROM golang:1.23 AS builder

ENV CGO_ENABLED=0

WORKDIR /builder

RUN mkdir server && mkdir stayalive

COPY ./server/main.go /builder/server

COPY /stayalive/ /builder/stayalive

COPY go.mod go.sum /builder/

RUN go build -tags netgo -a -v -o server ./server/

EXPOSE 2220 2222 2224

ENTRYPOINT [ "./server/server", "-p" ]

CMD [ "PORT" ]