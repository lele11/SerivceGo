# -*- coding: UTF-8 -*-
import socket
import os
import platform
import sys


try:
    import consul
except:
    os.system('pip install python-consul')
    import consul


consulAddr = "192.168.150.191"
port = "8500"
gameName = "dinosaur"


def isopen(addr, port):
    s = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
    try:
        s.connect((addr,port))
        s.shutdown(2)
        return True
    except:
        return False


def getip():
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        s.connect(('8.8.8.8', 80))
        ip = s.getsockname()[0]
    finally:
        s.close()
    return ip


def start(id, kind, tag):
    exe = sys.path[0] + "/main"
    system = platform.system()
    consuld = consulAddr + ":" + port
    if system == "Windows":
        os.system("start cmd /k %s.exe -ID %s -kind %s -consul %s %s " % (exe, id, kind, consuld, tag))
    else:
        os.system("nohup %s -ID %s -kind %s -consul %s %s  >>nohup.out_%s 2>&1 &" % (exe, id, kind, consuld, tag, tag))

if __name__ == '__main__':
    c = consul.Consul(host=consulAddr, port=port)
    data = c.agent.services()
    ip = getip()
    for k, v in data.items():
        if not v.get("Service", "").startswith(gameName):
            continue
        if v.get("Address", "") == ip:
            if isopen(ip, v.get("Port", 0)):
                print("server is open " + v.get("ID", ""))
            else:
                flag = v.get("ID", "")
                start(v["Meta"]["id"], v["Meta"]["kind"], flag)
    os.system("ps aux | grep %s | grep -v grep" % gameName)

