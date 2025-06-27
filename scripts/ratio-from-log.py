import csv
import os
import sys

here = os.getcwd()

#read from file
def countRatio(logFile):
    compHi = 0
    sentHi = 0
    compLo = 0
    sentLo = 0
    compBe = 0
    sentBe = 0
    print("Sorting data by TOS...")
    with open(logFile, mode='r') as data:
        dataFile = csv.reader(data, delimiter=' ')
        for line in dataFile:
            if "completed" in line and "hi" in line:
                with open("comp-hi.tmp", mode='a') as f:
                    f.write(str(line) + "\n")
            elif "completed" in line and "lo" in line:
                with open("comp-lo.tmp", mode='a') as f:
                    f.write(str(line) + "\n")
            elif "completed" in line and "be" in line:
                with open("comp-be.tmp", mode='a') as f:
                    f.write(str(line) + "\n")
            elif "priority:" in line and "hi" in line:
                with open("sent-hi.tmp", mode='a') as f:
                    f.write(str(line) + "\n")
            elif "priority:" in line and "lo" in line:
                with open("sent-lo.tmp", mode='a') as f:
                    f.write(str(line) + "\n")
            elif "priority:" in line and "be" in line:
                with open("sent-be.tmp", mode='a') as f:
                    f.write(str(line) + "\n")
    
    if os.path.isfile("comp-hi.tmp"):
        print("Calculating ratio from *-hi.tmp files...")
        with open("comp-hi.tmp", mode='r') as hiFile:
            compHi = len(hiFile.readlines())

    if os.path.isfile("sent-hi.tmp"):
        with open("sent-hi.tmp", mode='r') as hiFile:
            sentHi = len(hiFile.readlines())

    if os.path.isfile("comp-lo.tmp"):
        print("Calculating ratio from *-lo.tmp files...")
        with open("comp-lo.tmp", mode='r') as loFile:
            compLo = len(loFile.readlines())

        with open("sent-lo.tmp", mode='r') as loFile:
            sentLo = len(loFile.readlines())

    if os.path.isfile("comp-be.tmp"):
        print("Calculating ratio from *-be.tmp file...")
        with open("comp-be.tmp", mode='r') as beFile:
            compBe = len(beFile.readlines())

        with open("sent-be.tmp", mode='r') as beFile:
            sentBe = len(beFile.readlines())

    with open("ratio.csv", "w") as ratio_log:
        print("sent/comp,", "hi,", "lo,", "be", file=ratio_log)
        print("completed,", str(compHi)+',', str(compLo)+',', str(compBe), file=ratio_log)
        print("sent,", str(sentHi)+',', str(sentLo)+',', str(sentBe), file=ratio_log)
        # Avoid division by zero
        ratio_hi = compHi / sentHi if sentHi else 0
        ratio_lo = compLo / sentLo if sentLo else 0
        ratio_be = compBe / sentBe if sentBe else 0
        print("ratio,", str(ratio_hi)+',', str(ratio_lo)+',', str(ratio_be), file=ratio_log)

    for tmp_file in [
        "comp-hi.tmp", "sent-hi.tmp",
        "comp-lo.tmp", "sent-lo.tmp",
        "comp-be.tmp", "sent-be.tmp"
    ]:
        if os.path.isfile(tmp_file): os.remove(tmp_file)

    print("Ratio calculations completed. Logs saved to ratio.csv.")

if __name__ == "__main__":
    firstFilename = os.path.join(here, sys.argv[1])
    countRatio(firstFilename)


