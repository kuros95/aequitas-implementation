#!/bin/bash
file=$1
echo formatting $file...
while read -r line; do
    echo $line | tr " " , | tr - , >> $file.txt
done < $file