#!/bin/bash
#get output from tcpdump for a given socket and calculate it in mbps. length is in bytes

#gather data from the same second. Sum the bytes. Calculate bytes per second. 

file=$1

seconds=$(echo $(cat $file | cut -d "." -f 1 | awk '!seen[$0]++'))

for second in $seconds;
do
    bytes_in_second=$(echo $(cat $file | grep $second | cut -d " " -f 13 | cut -d ":" -f 1))
    bytes=0
    for i in $bytes_in_second; do
        bytes="$((bytes+i))"
    done
    echo $second, $bytes, Bps >> $file-stats.log
done
