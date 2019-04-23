# -*- coding: UTF-8 -*-
import common
import os
import platform


def start(configDir):
    if platform.system() == "Linux":
        os.system("nohup ./consul agent -config-file=%s -config-dir=%s &" % (common.configFile, configDir))
    if platform.system() == "Windows":
        os.system("cmd start /k consul.exe agent -config-file=%s -config-dir=%s" % (common.configFile, configDir))


if __name__ == '__main__':
    if not os.path.exists(common.configDir):
        exit("Add Config File use service.py")
    start(common.configDir)
