import os

processName = "dinosaur"


def stop():
    os.system("ps ux | grep %s  | grep -v grep | grep -v stop | grep -v publish |awk '{print $2}'| xargs  kill" %(processName))
    print("Stop Done")


if __name__ == '__main__':
    stop()
