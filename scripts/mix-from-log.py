import csv
import os
import sys

here = os.getcwd()

#read from file
def countRatio(logFile):
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

    # Collect all seconds in order, including those with zeroes
    all_seconds = set(mix.keys())
    # Also include seconds that appear in the log, even if no "completed" lines for them
    with open(logFile, mode='r') as data:
        dataFile = csv.reader(data, delimiter=' ')
        for line in dataFile:
            if len(line) > 1:
                all_seconds.add(line[1])

    with open('mix.csv', 'w') as f:
        headers = ["time", "hi", "lo", "all"]
        f.write(",".join(headers) + "\n")
        for second in sorted(all_seconds):
            hi = mix.get(second, {}).get("hi", 0)
            lo = mix.get(second, {}).get("lo", 0)
            f.write(f"{second},{hi},{lo},{hi+lo}\n")

    print("Mix calculations completed. Logs saved to mix.log.")

if __name__ == "__main__":
    firstFilename = os.path.join(here, sys.argv[1])
    countRatio(firstFilename)