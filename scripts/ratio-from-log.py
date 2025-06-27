import csv
import os
import sys

here = os.getcwd()

#read from file
def countRatio(logFile):

    compHi = compLo = compBe = 0
    sentHi = sentLo = sentBe = 0
    print("Sorting data by TOS...")

    with open(logFile, mode='r') as data:
        dataFile = csv.reader(data, delimiter=' ')
        for line in dataFile:
            # Use set for fast membership checks
            line_set = set(line)
            if "completed" in line_set:
                if "hi" in line_set:
                    compHi += 1
                elif "lo" in line_set:
                    compLo += 1
                elif "be" in line_set:
                    compBe += 1
            elif "priority:" in line_set:
                if "hi" in line_set:
                    sentHi += 1
                elif "lo" in line_set:
                    sentLo += 1
                elif "be" in line_set:
                    sentBe += 1

    with open("ratio.csv", "w") as ratio_log:
        print("sent/comp,", "hi,", "lo,", "be", file=ratio_log)
        print("completed,", f"{compHi},", f"{compLo},", f"{compBe}", file=ratio_log)
        print("sent,", f"{sentHi},", f"{sentLo},", f"{sentBe}", file=ratio_log)
        # Avoid division by zero
        ratio_hi = compHi / sentHi if sentHi else 0
        ratio_lo = compLo / sentLo if sentLo else 0
        ratio_be = compBe / sentBe if sentBe else 0
        print("ratio,", f"{ratio_hi},", f"{ratio_lo},", f"{ratio_be}", file=ratio_log)

    print("Ratio calculations completed. Logs saved to ratio.csv.")

if __name__ == "__main__":
    firstFilename = os.path.join(here, sys.argv[1])
    countRatio(firstFilename)


