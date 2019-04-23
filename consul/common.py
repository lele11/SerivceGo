import json
import os

configFile = "node.json"
configDir = "./config"


def readJsonFile(file):
    f = open(file, 'r')
    str = f.read()
    f.close()
    if str == "":
        return {}
    return json.loads(str)


def writeJson(file,data):
    d = json.dumps(data, sort_keys=True, indent=4, separators=(',', ': '))
    f = open(file, 'w+')
    f.write(d)
    f.close()


def initFile(file):
    dir = os.path.dirname(file)
    if not os.path.exists(dir):
        os.makedirs(dir)
    if not os.path.exists(file):
        f = open(file, 'w')
        f.close()


#游戏名称字典 GameName_SERVICE => KIND
serviceNameToKind = {
   
}
