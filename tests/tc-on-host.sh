#!/bin/bash
#Enable traffic classes on host. Cannot be done with DSCP, so using websockets.

src_ip=$(ip addr show eth0 scope global | grep 172 | awk '{print $2}' | cut -d '/' -f1)

echo "adding traffic control on eth0 interface and source ip" $src_ip

tc qdisc del dev eth0 root 2>/dev/null
tc qdisc add dev eth0 root handle 1: htb default 4 direct_qlen 100
tc class add dev eth0 parent 1: classid 1:1 htb rate 10mbps ceil 10mbps

tc class add dev eth0 parent 1:1 classid 1:2 htb rate 3mbps ceil 10mbps
tc class add dev eth0 parent 1:1 classid 1:3 htb rate 5mbps ceil 10mbps
tc class add dev eth0 parent 1:1 classid 1:4 htb rate 0mbps ceil 10mbps

tc filter add dev eth0 protocol ip parent 1: prio 1 u32 match ip dsfield 0x40 0x1e classid 1:2
tc filter add dev eth0 protocol ip parent 1: prio 2 u32 match ip dsfield 0x20 0x1e classid 1:3

tc qdisc add dev eth0 parent 1:2 handle 12: sfq perturb 10
tc qdisc add dev eth0 parent 1:3 handle 13: sfq perturb 10
tc qdisc add dev eth0 parent 1:4 handle 14: sfq perturb 10

echo "traffic control added"

# The int value for 0x20 is 32, and for 0x40 is 64.
#10 for 1, 10 for 10, 10 for 100, 10 for 1000
#for 30 it was 10/15/0 minimum
#for 35 it was 12/18/0 minimum
#for 40 it was 13/20/0 minimum
