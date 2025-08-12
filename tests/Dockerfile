FROM golang:1.23 AS tester

WORKDIR /tester

RUN apt-get update && apt-get install iproute2 iperf3 tcpdump vim -y 

COPY  tc-on-host.sh tc-test.sh /tester/

ENTRYPOINT [ "sleep", "infinity" ]
