FROM ubuntu:22.04 AS builder

WORKDIR /builder

COPY . /builder/
COPY ../stayalive /server/

RUN go build .

RUN ll /builder

FROM alpine:3.20 AS runner

WORKDIR /server

COPY builder/builder/main runner/server/

CMD [ "./main" ]