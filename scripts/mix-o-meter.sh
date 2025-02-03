#!/bin/bash

file=$1
echo time, hi, lo, all > $file-mix.log
seconds=$(echo $(cat $file | cut -d " " -f 2 | awk '!seen[$0]++'))

for second in $seconds; do
    mix_in_second=$(echo $(cat $file | grep $second | grep completed | cut -d " " -f 11))
    hi=0
    lo=0
    for i in $mix_in_second; do
        if [ $i = 'hi' ]; then
            hi="$((hi+1))"
        elif [ $i = 'lo' ]; then
            lo="$((lo+1))"
        fi
    done
    all="$((hi+lo))"
    echo $second, $hi, $lo, $all >> $file-mix.log
done
