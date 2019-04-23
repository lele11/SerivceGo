package db

/*
	accountKey:用来保存用户账号相关的数据，玩家账号应相对独立于玩家的游戏数据
	key格式： account:channel:accountName
	类型： hash
	Field   Value
	SID   	UID					// 区服id 对应 玩家的uid  可以有多组
	TYPE 	1					// 账号的类型，可以区分管理员账号，普通账号
	RECENT  [SID1,SID2...]		// 最近登录的区服列表，有序数组
*/

import (
	"game/base/redis"
)

const (
	AccountKey         = "account"
	AccountFieldRecent = "recent"
	AccountFieldType   = "type"
)

func getAccountKey(channel, account string) string {
	return AccountKey + ":" + channel + ":" + account
}

// GetAccountUID 账号数据结构  hash  keys ACCOUNT:CHANNEL  field ServerID value uid
func GetAccountUID(loginChannel string, account string, serverID uint64) (uint64, bool) {
	key := getAccountKey(loginChannel, account)
	id := redis.HGetUint64(key, serverID)
	if id > 0 {
		return id, false
	}
	id = uint64(GetUniqueId("useruid") + 100000)
	//TODO 单服id自增方式
	// id = 0
	redis.HSet(key, serverID, id)

	return uint64(id), true
}

// SetAccountType 设置账号类型
func SetAccountType(channel, account string, aType uint32) {
	redis.HSet(getAccountKey(channel, account), AccountFieldType, aType)
}

// AddAccountRecent 账号最近登录的区服
func AddAccountRecent(channel, account string, serverID uint64) {
	key := getAccountKey(channel, account)
	var r []uint64
	redis.HGetToStruct(key, AccountFieldRecent, &r)
	for k, v := range r {
		if v == serverID {
			r = append(r[:k], r[k+1:]...)
			break
		}
	}
	r = append(r, serverID)
	redis.HSetToJson(key, AccountFieldRecent, r)
}

// GetAccountRecent 获取最近区服列表
func GetAccountRecent(channel, account string) (r []uint64) {
	redis.HGetToStruct(getAccountKey(channel, account), AccountFieldRecent, &r)
	return
}

// GetUniqueId 获取唯一ID
func GetUniqueId(kind string) int64 {
	return redis.HIncrby("UniqueId", kind, 1)
}
