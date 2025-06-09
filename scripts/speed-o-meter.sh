#!/bin/bash
#get output from tcpdump for a given socket and calculate it in mbps. length is in bytes

#gather data from the same second. Sum the bytes. Calculate bytes per second. 

file=$1

cat $file | grep "tos 0x40" > hi-log.tmp
cat $file | grep "tos 0x20" > lo-log.tmp
cat $file | grep "tos 0x10" > be-log.tmp

echo calculating speed from hi-log.tmp in $file...
seconds=$(echo $(cat hi-log.tmp | cut -d "." -f 1 | awk '!seen[$0]++'))
for second in $seconds;
do
    bytes_in_second=$(echo $(cat hi-log.tmp | grep $second | cut -d " " -f 13 | cut -d ":" -f 1))
    bytes=0
    for i in $bytes_in_second; do
        bytes="$((bytes+i))"
    done
    echo $second, $bytes, Bps >> $file-hi-log-speed.log
done

echo calculating speed from lo-log.tmp in $file...
seconds=$(echo $(cat lo-log.tmp | cut -d "." -f 1 | awk '!seen[$0]++'))
for second in $seconds;
do
    bytes_in_second=$(echo $(cat lo-log.tmp | grep $second | cut -d " " -f 13 | cut -d ":" -f 1))
    bytes=0
    for i in $bytes_in_second; do
        bytes="$((bytes+i))"
    done
    echo $second, $bytes, Bps >> $file-lo-log-speed.log
done

echo calculating speed from be-log.tmp in $file...
seconds=$(echo $(cat be-log.tmp | cut -d "." -f 1 | awk '!seen[$0]++'))
for second in $seconds;
do
    bytes_in_second=$(echo $(cat be-log.tmp | grep $second | cut -d " " -f 13 | cut -d ":" -f 1))
    bytes=0
    for i in $bytes_in_second; do
        bytes="$((bytes+i))"
    done
    echo $second, $bytes, Bps >> $file-be-log-speed.log
done

rm hi-log.tmp lo-log.tmp be-log.tmp