package db

import (
	"fmt"
	"game/base/redis"
	"game/base/util"
	"time"
)

func sessionSet(uid uint64) string {
	return fmt.Sprintf("session:%d", uid)
}
func StoreSessionInfo(uid uint64, sessionKey string, sid uint64) {
	redis.HMSet(sessionSet(uid), map[interface{}]interface{}{
		"sessionKey":    sessionKey,
		"sessionExpire": time.Now().Unix() + 1800,
		"sid":           sid,
	})
}
func GetSessionInfo(uid uint64) (sessionKey string, sessionExpire int64, sid uint64) {
	d := redis.HMGetSlice(sessionSet(uid), []interface{}{"sessionKey", "sessionExpire", "sid"})
	if len(d) != 3 {
		return
	}
	sessionKey = d[0]
	sessionExpire = util.StringToInt64(d[1])
	sid = util.StringToUint64(d[2])
	return
}

func SetGateWaySrv(uid uint64, sid uint64) {
	redis.HSet(sessionSet(uid), "gatewaySrv", sid)
}

func GetGateWaySrv(uid uint64) uint64 {
	return redis.HGetUint64(sessionSet(uid), "gatewaySrv")
}

func SetGameSrv(uid uint64, value interface{}) {
	redis.HSet(sessionSet(uid), "gameSrv", value)
}
