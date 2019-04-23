package redis

import (
	"encoding/json"
	"errors"
	"game/base/util"

	"reflect"

	"github.com/cihub/seelog"
)

var ErrNotFound = errors.New("NotFoundValue")

// Exists key是否存在
func Exists(key string) bool {
	return exists(key)
}

// Del 删除key
func Del(key string) {
	delKey(key)
}

// HVals hash 值
func HVals(key string) []string {
	return hvals(key)
}

// HKeysUint64 hash键 uint64 slice
func HKeysUint64(key string) []uint64 {
	var r []uint64
	for _, v := range hKeys(key) {
		r = append(r, util.StringToUint64(v))
	}
	return r
}

// HKeysUint32 hash键 uint32 slice
func HKeysUint32(key string) []uint32 {
	var r []uint32
	for _, v := range hKeys(key) {
		r = append(r, util.StringToUint32(v))
	}
	return r
}

// HGetFloat64 hash value float64
func HGetFloat64(key string, field interface{}) float64 {
	value := hGet(key, field)
	return util.StringToFloat64(value)
}

// HGetUint64 hash value uint64
func HGetUint64(key string, field interface{}) uint64 {
	value := hGet(key, field)
	return util.StringToUint64(value)
}

// HGetUint32 hash value uint32
func HGetUint32(key string, field interface{}) uint32 {
	value := hGet(key, field)
	return util.StringToUint32(value)
}

// HGetInt64 hash value int64
func HGetInt64(key string, field string) int64 {
	value := hGet(key, field)
	return util.StringToInt64(value)
}

// HGetInt hash value int
func HGetInt(key string, field interface{}) int {
	value := hGet(key, field)
	return util.StringToInt(value)
}

// HGetUint16 hash value uint16
func HGetUint16(key string, field string) uint16 {
	value := hGet(key, field)
	return util.StringToUint16(value)
}

// HMGetMap hash value map[string]string
func HMGetMap(key string, field []string) map[string]string {
	var ff []interface{}
	for _, f := range field {
		ff = append(ff, f)
	}
	r := hMGet(key, ff)
	m := map[string]string{}
	if len(r) == 0 {
		return m
	}
	for k, v := range field {
		m[v] = r[k]
	}
	return m
}

// HMGetSlice hash value []string
func HMGetSlice(key string, field []interface{}) []string {
	return hMGet(key, field)
}

// HIncrby
func HIncrby(key string, field interface{}, value interface{}) int64 {
	return hIncrBy(key, field, value)
}

// HSetToJson 对象保存为json结构
func HSetToJson(key string, field, value interface{}) {
	d, e := json.Marshal(value)
	if e != nil {
		seelog.Errorf("HSet TO Json key %s filed %s value %v Error %s", key, field, value, e)
		return
	}
	hSet(key, field, d)
}

// HGetToStruct 获取hash值 解析json
func HGetToStruct(key string, field interface{}, dst interface{}) error {
	value := hGet(key, field)
	if value == "" {
		return ErrNotFound
	}
	e := json.Unmarshal([]byte(value), dst)
	if e != nil {
		seelog.Errorf("Hget To Struct Error  key %s field %s error ", key, field, e)
	}
	return e
}

// HGetAll map[string]string
func HGetAll(key string) map[string]string {
	return hGetAll(key)
}

// HGetAllMapU32U64 hGetAll map[uint32]uint64
func HGetAllMapU32U64(key string) map[uint32]uint64 {
	m := make(map[uint32]uint64)
	for k, v := range hGetAll(key) {
		m[util.StringToUint32(k)] = util.StringToUint64(v)
	}
	return m
}

// HGetAllMapU32U32 hGetAll map[uint32]uint32
func HGetAllMapU32U32(key string) map[uint32]uint32 {
	m := make(map[uint32]uint32)
	for k, v := range hGetAll(key) {
		m[util.StringToUint32(k)] = util.StringToUint32(v)
	}
	return m
}

// HGetAllMapU64U32 hGetAll map[uint64]uint32
func HGetAllMapU64U32(key string) map[uint64]uint32 {
	m := make(map[uint64]uint32)
	for k, v := range hGetAll(key) {
		m[util.StringToUint64(k)] = util.StringToUint32(v)
	}
	return m
}

// HGetAllToStruct hGetAll []interface{}
func HGetAllToStruct(key string, dst interface{}) []interface{} {
	value := hGetAll(key)
	dstType := reflect.TypeOf(dst).Elem()
	var ret []interface{}
	for _, str := range value {
		if str[0] != 123 {
			continue
		}
		m := reflect.New(dstType).Interface()
		e := json.Unmarshal([]byte(str), m)
		if e != nil {
			seelog.Errorf("HgetAll To Struct Error  key %s field %s error ", key, e)
			continue
		}
		ret = append(ret, m)

	}
	return ret
}

// HGet to string
func HGet(key string, field interface{}) string {
	return hGet(key, field)
}

// HSet hash 设置
func HSet(key string, field interface{}, value interface{}) {
	hSet(key, field, value)
}

// HMSet HMSet
func HMSet(key string, field map[interface{}]interface{}) {
	hMSet(key, field)
}

// HDEL 删除
func HDEL(key string, field interface{}) {
	hDEL(key, field)
}

// HExists field是否存在
func HExists(key string, field interface{}) bool {
	return hExists(key, field)
}

// HLen hash长度
func HLen(key string) uint64 {
	return hLen(key)
}

// HSetNx 如果没有则设置成功
func HSetNx(key string, field, value interface{}) bool {
	return hSetNx(key, field, value)
}

// ZRevRank 获取排名
func ZRevRank(key string, member interface{}) uint64 {
	return zRevRank(key, member)
}

// ZRevRange 获取排行信息 ，从大到小
func ZRevRange(key string, start, end interface{}, isScore bool) []string {
	return zRevRange(key, start, end, isScore)
}

// ZAdd 添加zset
func ZAdd(key string, field, value interface{}) {
	zAdd(key, field, value)
}

// ZCard 集合大小
func ZCard(key string) uint64 {
	return zCard(key)
}

// ZRangeByScore 按分值排名
func ZRangeByScore(key string, start, end interface{}) []string {
	return zRangeByScore(key, start, end, false)
}

// ZRangeByScoreWithLimit 设置获取列表大小
func ZRangeByScoreWithLimit(key string, start, end, limit interface{}) []string {
	return zRangeByScoreWithLimit(key, start, end, limit)
}

// ZScore
func ZScore(key string, member interface{}) string {
	return zScore(key, member)
}

// ExpireKey
func ExpireKey(key string, dur int64) {
	expireKey(key, dur)
}

// ExpireKeyAt
func ExpireKeyAt(key string, dur int64) {
	expireKeyAt(key, dur)

}

// LPush
func LPush(key string, value interface{}) {
	lPush(key, value)
}

// LPop
func LPop(key string) string {
	return lPop(key)
}

// LRange
func LRange(key string, start, end int32) []string {
	return lRange(key, start, end)
}

// LTRIM
func LTRIM(key string, start, end int32) {
	lTRIM(key, start, end)
}

// SAdd
func SAdd(key string, member interface{}) {
	sAdd(key, member)
}

// SMembers
func SMembers(key string) []string {
	return sMembers(key)
}

// SIsMember
func SIsMember(key string, member interface{}) bool {
	return sIsMember(key, member)
}

// SRem
func SRem(key string, member interface{}) {
	sRem(key, member)
}

// SRandMember
func SRandMember(key string, count interface{}) []string {
	return sRandMember(key, count)
}
