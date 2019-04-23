import argparse
import os

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-f', '--config', required=True,  help=u"(必须)配置文件")
    args = parser.parse_args()

    file = open(args.config, "r")
    for cmd in file.readlines():
        if cmd == "":
            continue
        p = os.system(cmd)
        print("run %s result %s" % (cmd, p))

