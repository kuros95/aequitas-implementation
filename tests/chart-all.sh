#!/bin/bash

foreach file $(ls results) do
    ./get-results.sh "results/$file"
done

# foreach newfile $(ls results | grep '\.txt$' ) do

#     python3 chart-data.py 
# done
