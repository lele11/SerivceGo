package redis

import (
	"encoding/json"
	"errors"
	"game/base/util"

	"reflect"

	"github.com/cihub/seelog"
)

var ErrNotFound = errors.New("NotFoundValue")

func Exists(key string) bool {
	return exists(key)
}
func Del(key string) {
	delKey(key)
}

func HVals(key string) []string {
	return hvals(key)
}

func HKeysUint64(key string) []uint64 {
	var r []uint64
	for _, v := range hKeys(key) {
		r = append(r, util.StringToUint64(v))
	}
	return r
}
func HKeysUint32(key string) []uint32 {
	var r []uint32
	for _, v := range hKeys(key) {
		r = append(r, util.StringToUint32(v))
	}
	return r
}

// Hash
func HGetFloat64(key string, field interface{}) float64 {
	value := hGet(key, field)
	return util.StringToFloat64(value)
}
func HGetUint64(key string, field interface{}) uint64 {
	value := hGet(key, field)
	return util.StringToUint64(value)
}

func HGetUint32(key string, field interface{}) uint32 {
	value := hGet(key, field)
	return util.StringToUint32(value)
}

func HGetInt64(key string, field string) int64 {
	value := hGet(key, field)
	return util.StringToInt64(value)
}

func HGetInt(key string, field interface{}) int {
	value := hGet(key, field)
	return util.StringToInt(value)
}

func HGetUint16(key string, field string) uint16 {
	value := hGet(key, field)
	return util.StringToUint16(value)
}

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
func HMGetSlice(key string, field []interface{}) []string {
	return hMGet(key, field)
}

func HIncrby(key string, field interface{}, value interface{}) int64 {
	return hIncrBy(key, field, value)
}

func HSetToJson(key string, field, value interface{}) {
	d, e := json.Marshal(value)
	if e != nil {
		seelog.Errorf("HSet TO Json key %s filed %s value %v Error %s", key, field, value, e)
		return
	}
	hSet(key, field, d)
}

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

func HGetAll(key string) map[string]string {
	return hGetAll(key)
}
func HGetAllMapU32U64(key string) map[uint32]uint64 {
	m := make(map[uint32]uint64)
	for k, v := range hGetAll(key) {
		m[util.StringToUint32(k)] = util.StringToUint64(v)
	}
	return m
}
func HGetAllMapU32U32(key string) map[uint32]uint32 {
	m := make(map[uint32]uint32)
	for k, v := range hGetAll(key) {
		m[util.StringToUint32(k)] = util.StringToUint32(v)
	}
	return m
}
func HGetAllMapU64U32(key string) map[uint64]uint32 {
	m := make(map[uint64]uint32)
	for k, v := range hGetAll(key) {
		m[util.StringToUint64(k)] = util.StringToUint32(v)
	}
	return m
}
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

func HGet(key string, field interface{}) string {
	return hGet(key, field)
}

func HSet(key string, field interface{}, value interface{}) {
	hSet(key, field, value)
}
func HMSet(key string, field map[interface{}]interface{}) {
	hMSet(key, field)
}

func HDEL(key string, field interface{}) {
	hDEL(key, field)
}
func HExists(key string, field interface{}) bool {
	return hExists(key, field)
}
func HLen(key string) uint64 {
	return hLen(key)
}

func HSetNx(key string, field, value interface{}) bool {
	return hSetNx(key, field, value)
}

// sorted set
func ZRevRank(key string, member interface{}) uint64 {
	return zRevRank(key, member)
}

func ZRevRange(key string, start, end interface{}, isScore bool) []string {
	return zRevRange(key, start, end, isScore)
}
func ZAdd(key string, field, value interface{}) {
	zAdd(key, field, value)
}
func ZCard(key string) uint64 {
	return zCard(key)
}
func ZRangeByScore(key string, start, end interface{}) []string {
	return zRangeByScore(key, start, end, false)
}
func ZRangeByScoreWithLimit(key string, start, end, limit interface{}) []string {
	return zRangeByScoreWithLimit(key, start, end, limit)
}
func ZScore(key string, member interface{}) string {
	return zScore(key, member)
}
func ExpireKey(key string, dur int64) {
	expireKey(key, dur)
}
func ExpireKeyAt(key string, dur int64) {
	expireKeyAt(key, dur)

}

// List
func LPush(key string, value interface{}) {
	lPush(key, value)
}

func LPop(key string) string {
	return lPop(key)
}

func LRange(key string, start, end int32) []string {
	return lRange(key, start, end)
}

func LTRIM(key string, start, end int32) {
	lTRIM(key, start, end)
}

//==========================

func SAdd(key string, member interface{}) {
	sAdd(key, member)
}
func SMembers(key string) []string {
	return sMembers(key)
}
func SIsMember(key string, member interface{}) bool {
	return sIsMember(key, member)
}
func SRem(key string, member interface{}) {
	sRem(key, member)
}
func SRandMember(key string, count interface{}) []string {
	return sRandMember(key, count)
}
