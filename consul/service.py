# -*- coding: UTF-8 -*-
import common
import random
import argparse
import consul


def getPort(ports):
    done = False
    p = 0
    while not done:
        p = random.randint(6000, 7000)
        e = False
        for i in ports:
            if i == p:
                e = True
                break
        if not e:
            done = True
    return p


def update(configFile, kind, id, meta={}):
    fileInfo = common.readJsonFile(configFile)
    if "services" not in fileInfo:
        return

    Id = kind + "_" + id
    sInfo = {}
    for val in fileInfo["services"]:
        if val["id"] == Id:
            sInfo = val
            fileInfo["services"].remove(val)
            break
    sInfo["meta"].update(meta)
    fileInfo['services'].append(sInfo)
    common.writeJson(configFile, fileInfo)


def remove(configFile, kind, id):
    fileInfo = common.readJsonFile(configFile)
    if "services" not in fileInfo:
        return

    Id = kind + "_" + id
    for val in fileInfo["services"]:
        if val["id"] == Id:
            fileInfo["services"].remove(val)
            break
    common.writeJson(configFile, fileInfo)
    c = consul.Consul(port=8500)
    c.agent.service.deregister(Id)


def create(configFile, kind, id, type, ip, port = 0, meta = {}):
    common.initFile(configFile)
    fileInfo = common.readJsonFile(configFile)
    if "services" not in fileInfo:
        fileInfo["services"] = []

    Id = kind + "_" + id
    ports = []
    for val in fileInfo["services"]:
        if val["id"] == Id:
            return
        if val["address"] == ip:
            ports.append(val["port"])
            ports.append(val["meta"]["id"])
    if port == 0:
        port = getPort(ports)

    s = {
        "id": Id,
        "name": kind,
        "meta": {
                "id": id,
                "kind": type
        },
        "address": ip,
        "port": int(port),
        "check": {
            "id": Id,
            "name": Id,
            "ttl": "10s"
        }
    }
    s["meta"].update(meta)
    fileInfo['services'].append(s)
    common.writeJson(configFile, fileInfo)


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-d', '--do', required=True, choices=["create", "update", "delete"], help=u"(必须)操作行为")
    parser.add_argument('-k', '--kind', required=True, choices=common.serviceNameToKind.keys(), help=u"(必须)服务的类型： Game ，由游戏配置")
    parser.add_argument('-i', '--id', required=True,  default="", help=u"(必须)服务的id：3001 ，全局唯一")
    parser.add_argument('-a', '--address', default="", help=u"(创建必须)服务的地址")
    parser.add_argument('-p', '--port', default=0, help=u"(可选)服务的端口，不填则随机端口")
    parser.add_argument('-m', '--meta', nargs='*', default="", help=u" (可选)服务的其他属性")
    args = parser.parse_args()
    # TODO 检查ip格式 端口范围
    meta = {}
    for v in args.meta:
        i = v.split('=')
        meta[i[0]] = i[1]

    file, _ = args.kind.split("_")
    kind = args.kind
    configFile = common.configDir + "/" + file + ".json"
    if args.do == "create":
        if args.address == "":
            exit('need address ')
        create(configFile, kind, args.id, common.serviceNameToKind[kind], args.address, args.port, meta)
    if args.do == "update":
        update(configFile, kind, args.id, meta)
    if args.do == "delete":
        remove(configFile, kind, args.id)

    # serviceTemp = '''{
    #     	"services":[
    #     	 {
    #     		"id": "Battle_5001",
    #     		"name": "Battle",
    #     		"meta": {
    #     			"id": "5001",
    #     			"kind": "5"
    #     		},
    #     		"address": "192.168.92.75",
    #     		"port": 0,
    #     		"check": {
    #     			"id": "Battle_5001",
    #     			"name": "Battle_5001",
    #     			"ttl": "10s"
    #     		}
    #     	}
    #     	]
    #     }
    # '''