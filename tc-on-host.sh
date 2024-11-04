#!/bin/sh
#Enable traffic classes on host. Cannot be done with DSCP, so using websockets.


tc qdisc add dev eth0 root handle 1: htb default 4 direct_qlen 100000
tc class add dev eth0 parent 1: classid 1:1 htb rate 1000mbps ceil 1000mbps

tc class add dev eth0 parent 1:1 classid 1:2 htb rate 200mbps ceil 1000mbps
tc class add dev eth0 parent 1:1 classid 1:3 htb rate 700mbps ceil 1000mbps
tc class add dev eth0 parent 1:1 classid 1:4 htb rate 100mbps ceil 1000mbps

tc filter add dev eth0 protocol ip parent 1: prio 1 u32 match ip src $src_ip \
        match ip dport 2220 0xffff flowid 1:2
tc filter add dev eth0 protocol ip parent 1: prio 2 u32 match ip src $src_ip \
        match ip dport 2222 0xffff flowid 1:3

tc qdisc add dev eth0 parent 1:2 handle 12: sfq perturb 10
tc qdisc add dev eth0 parent 1:3 handle 13: sfq perturb 10
tc qdisc add dev eth0 parent 1:4 handle 14: sfq perturb 10

