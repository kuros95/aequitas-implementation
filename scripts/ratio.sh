#!/bin/bash
#This script counts the number of completed and sent tasks for each priority level (hi, lo, be) in a given file.
file=$1


cat $file | grep 'completed' | grep 'hi' > hi-c.tmp

cat $file | grep 'completed' | grep 'lo' > lo-c.tmp

cat $file | grep 'completed' | grep 'be' > be-c.tmp

compHi=$(wc -l < hi-c.tmp)

compLo=$(wc -l < lo-c.tmp)

compBe=$(wc -l < be-c.tmp)


cat $file | grep 'priority' | grep 'hi' > hi-s.tmp

cat $file | grep 'priority' | grep 'lo' > lo-s.tmp

cat $file | grep 'priority' | grep 'be' > be-s.tmp

sentHi=$(wc -l < hi-s.tmp)

sentLo=$(wc -l < lo-s.tmp)

sentBe=$(wc -l < be-s.tmp)


echo s/c,hi,lo,be > $file-count.log
echo c,$compHi,$compLo,$compBe >> $file-count.log
echo s,$sentHi,$sentLo,$sentBe >> $file-count.log

rm hi-c.tmp hi-s.tmp lo-c.tmp lo-s.tmp be-c.tmp be-s.tmp