import matplotlib.pyplot as plt
import numpy as np
import csv
import os
import sys

here = os.path.dirname(os.path.abspath(__file__))

#read from file
speeds = list()
secondSpeeds = list()
seconds = list()
secondSeconds = list()
def readAndChart(firstFile, secondFile):
    with open(firstFile, mode = 'r') as data:
        dataFile = csv.reader(data)
        for line in dataFile:
            if "MBytes" in line:
                second = line[2]
                speed = line[5]
                seconds.append(float(second))
                speeds.append(float(speed))

    with open(secondFile, mode = 'r') as data:
        dataFile = csv.reader(data)
        for line in dataFile:
            if "MBytes" in line:
                second = line[2]
                speed = line[5]
                secondSeconds.append(float(second))
                secondSpeeds.append(float(speed))


    x = tuple(seconds)
    y = tuple(speeds)
    x1 = tuple(secondSeconds)
    y1 = tuple(secondSpeeds)
    plt.plot(x, y)
    plt.plot(x1, y1)
    plt.xlabel("Czas [s]")
    plt.ylabel("Przepływność [Mb/s]")
    plt.xticks(np.arange(0, 60, step=10))
    plt.yticks(np.arange(0, 100, step=10))
    # ax = plt.gca()
    # ax.set_ylim([0, 100])
    plt.show()
    plt.savefig(firstFile + ".jpg")

if __name__ == "__main__":
    firstFilename = os.path.join(here, sys.argv[1])
    secondFilename = os.path.join(here, sys.argv[2])
    readAndChart(firstFilename, secondFilename)


