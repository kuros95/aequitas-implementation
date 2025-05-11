#!/bin/bash
#Enable traffic classes on host. Cannot be done with DSCP, so using websockets.

src_ip=$(ip addr show eth0 scope global | grep 172 | awk '{print $2}')

echo "adding traffic control on eth0 interface and source ip" $src_ip

tc qdisc add dev eth0 root handle 1: htb default 2 direct_qlen 100
tc class add dev eth0 parent 1: classid 1:1 htb rate 100mbps ceil 100mbps

tc class add dev eth0 parent 1:1 classid 1:2 htb rate 1mbps ceil 100mbps
tc class add dev eth0 parent 1:1 classid 1:3 htb rate 1mbps ceil 100mbps

tc filter add dev eth0 protocol ip parent 1: u32 match ip dsfield 0x10 0x1e classid 1:2
tc filter add dev eth0 protocol ip parent 1: u32 match ip dsfield 0x8 0x1e classid 1:3

tc qdisc add dev eth0 parent 1:2 handle 12: sfq perturb 10
tc qdisc add dev eth0 parent 1:3 handle 13: sfq perturb 10

echo "traffic control added"

#iperf3 -t 60 -c 172.17.0.3 -p 2222 --logfile b-5-5-2
#iperf3 -s -p 2222
#iperf3 -t 60 -c 172.17.0.3 -p 2222 -u -b 4G --logfile b-5-5-2-u