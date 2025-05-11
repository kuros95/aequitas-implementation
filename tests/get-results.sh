#!/bin/bash
file=$1
echo calculating speed from $file...
while read -r line; do
    echo $line | tr " " , | tr - , >> $file.txt
done < $file