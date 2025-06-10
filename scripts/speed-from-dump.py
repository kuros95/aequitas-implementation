import csv
import os
import sys

here = os.getcwd()

#read from file
def sortByTos(dumpfile):
    print("Sorting data by TOS...")
    hiTmp = list()
    loTmp = list()
    beTmp = list()
    with open(dumpfile, mode='r') as data:
        dataFile = csv.reader(data, delimiter=' ')
        for line in dataFile:
            if "172.17.0.2.2222" in line:
                hiTmp.append(line)
            elif "172.17.0.2.2223" in line:
                loTmp.append(line)
            elif "172.17.0.2.2224" in line:
                beTmp.append(line)

    if len(hiTmp) > 0:
        print("Calculating speed from data-hi.tmp file...")
        seconds = list()
        for line in hiTmp:
            if line[0].split('.')[0] not in seconds:
                seconds.append(line[0].split('.')[0])
        
        bytesInSecond = 0
        with open("speed-hi.csv", mode='w') as speedHiFile:
            print("time,", "speed-in-Bytes", file=speedHiFile)
        for second in seconds:
            for line in hiTmp:
                if second in line[0] and isinstance(int(line[12].split(':')[0]), int):
                    bytes = line[12].split(':')[0]
                    bytesInSecond += int(bytes)
            with open("speed-hi.csv", mode='a') as speedHiFile:
                print(second+',', bytesInSecond, file=speedHiFile)
            bytesInSecond = 0

    if len(loTmp) > 0:
        print("Calculating speed from data-lo.tmp file...")
        seconds = list()
        for line in loTmp:
            if line[0].split('.')[0] not in seconds:
                seconds.append(line[0].split('.')[0])
        
        bytesInSecond = 0
        with open("speed-lo.csv", mode='w') as speedLoFile:
            print("time,", "speed-in-Bytes", file=speedLoFile)
        for second in seconds:
            for line in loTmp:
                if second == line[0].split('.')[0] and isinstance(int(line[12].split(':')[0]), int):
                    bytes = line[12].split(':')[0]
                    bytesInSecond += int(bytes)
            with open("speed-log.csv", mode='a') as speedLoFile:
                print(second+',', bytesInSecond, file=speedLoFile)
            bytesInSecond = 0

    if len(beTmp) > 0:
        print("Calculating speed from data-be.tmp file...")
        seconds = list()
        for line in beTmp:
            if line[0].split('.')[0] not in seconds:
                seconds.append(line[0].split('.')[0])
        
        bytesInSecond = 0
        with open("speed-be.csv", mode='w') as speedBeFile:
            print("time,", "speed-in-Bytes", file=speedBeFile)
        for second in seconds:
            for line in beTmp:
                if second == line[0].split('.')[0] and isinstance(int(line[12].split(':')[0]), int):
                    bytes = line[12].split(':')[0]
                    bytesInSecond += int(bytes)
            with open("speed-be.csv", mode='a') as speedBeFile:
                print(second+',', bytesInSecond, file=speedBeFile)
            bytesInSecond = 0
    
    print("Speed calculations completed. Logs saved to speed-hi.csv, speed-lo.csv, and speed-be.csv.")

if __name__ == "__main__":
    firstFilename = os.path.join(here, sys.argv[1])
    sortByTos(firstFilename)


