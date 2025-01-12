#!/bin/bash

file=$1

cat $file | grep 'completed' > completed.tmp

cat completed.tmp | grep 'hi' > hi-c.tmp

cat completed.tmp | grep 'lo' > lo-c.tmp

cat completed.tmp | grep 'be' > be-c.tmp

compHi=$(wc -l < hi-c.tmp)

compLo=$(wc -l < lo-c.tmp)

compBe=$(wc -l < be-c.tmp)


cat $file | grep 'priority' > sent.tmp

cat sent.tmp | grep 'hi' > hi-s.tmp

cat sent.tmp | grep 'lo' > lo-s.tmp

cat sent.tmp | grep 'be' > be-s.tmp

sentHi=$(wc -l < hi-s.tmp)

sentLo=$(wc -l < lo-s.tmp)

sentBe=$(wc -l < be-s.tmp)


echo s/c,hi,lo,be > $file-count.log
echo c,$compHi,$compLo,$compBe >> $file-count.log
echo s,$sentHi,$sentLo,$sentBe >> $file-count.log

rm completed.tmp sent.tmp hi-c.tmp hi-s.tmp lo-c.tmp lo-s.tmp be-c.tmp be-s.tmp