import csv
import os
import sys

here = os.getcwd()

#read from file
def sortByTos(dumpfile):
    print("Catching time...")
    seconds = set()

    with open("client.log", mode='r') as data:
        dataFile = csv.reader(data, delimiter=' ')
        for line in dataFile:
            if line[1] not in seconds:
                seconds.add(line[1])   

    print("Calculating speed...")

    # Prepare counters for each second
    speed = {sec: {'hi': 0, 'lo': 0, 'be': 0} for sec in seconds}

    with open(dumpfile, mode='r') as data:
        dataFile = csv.reader(data, delimiter=' ')
        for line in dataFile:
            try:
                sec = line[0].split('.')[0]
                bytes_val = int(line[12].split(':')[0])
            except (IndexError, ValueError):
                continue
            if sec in speed:
                if "172.17.0.2.2222" in line:
                    speed[sec]['hi'] += bytes_val
                elif "172.17.0.2.2223" in line:
                    speed[sec]['lo'] += bytes_val
                elif "172.17.0.2.2224" in line:
                    speed[sec]['be'] += bytes_val

    with open("speed.csv", mode='w') as speedFile:
        print("time,speed-hi,speed-lo,speed-be,sum", file=speedFile)
        for sec in sorted(speed.keys()):
            hi = speed[sec]['hi']
            lo = speed[sec]['lo']
            be = speed[sec]['be']
            print(f"{sec},{hi},{lo},{be},{hi+lo+be}", file=speedFile)
    
    print("Speed calculations completed. Logs saved to speed.csv.")

if __name__ == "__main__":
    firstFilename = os.path.join(here, sys.argv[1])
    sortByTos(firstFilename)


