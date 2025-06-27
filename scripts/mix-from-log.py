import csv
import os
import sys

here = os.getcwd()

#read from file
def countRatio(logFile):
    # seconds = list()
    # lines = list()
    print("Calculating mix from client.log...")

    mix = {}

    with open(logFile, mode='r') as data:
        dataFile = csv.reader(data, delimiter=' ')
        for line in dataFile:
            if len(line) > 11 and line[2] == "completed":
                second = line[1]
                if second not in mix:
                    mix[second] = {"hi": 0, "lo": 0}
                if line[10] == "hi":
                    mix[second]["hi"] += 1
                elif line[10] == "lo":
                    mix[second]["lo"] += 1

    with open('mix.csv', 'w') as f:
        headers = ["time", "hi", "lo", "all"]
        f.write(",".join(headers) + "\n")
        for second in sorted(mix.keys()):
            hi = mix[second]["hi"]
            lo = mix[second]["lo"]
            f.write(f"{second},{hi},{lo},{hi+lo}\n")


    # with open(logFile, mode='r') as data:
    #     dataFile = csv.reader(data, delimiter=' ')
    #     for line in dataFile:
    #         lines.append(line)
    #         if line[1] not in seconds:
    #             seconds.append(line[1])
    #     with open('mix.csv', 'w') as f:
    #         headers = ["time", "hi", "lo", "all"]
    #         f.write(",".join(headers) + "\n")

    # for second in seconds:
    #     hi = 0
    #     lo = 0
    #     for line in lines:
    #         if len(line) > 11:
    #             if line[1] == second and line[2] == "completed" and line[10] == "hi":
    #                 hi+=1
    #             elif line[1] == second and line[2] == "completed" and line[10] == "lo":
    #                 lo+=1
    #     with open('mix.csv', 'a') as f:
    #         row = [second, str(hi), str(lo), str(hi + lo)]
    #         f.write(",".join(row) + "\n")

    print("Mix calculations completed. Logs saved to mix.log.")

if __name__ == "__main__":
    firstFilename = os.path.join(here, sys.argv[1])
    countRatio(firstFilename)


