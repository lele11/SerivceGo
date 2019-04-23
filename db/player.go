package db

import (
	"game/base/redis"
	"game/base/util"
	"game/protoMsg"
)

const (
	ChannelField  = "channel"  //
	NickNameField = "nickName" //
	PictureField  = "picture"  //
	GenderField   = "gender"
	CountryField  = "country"
	ProvinceField = "province" //
	CityField     = "city"     //
	LanguageField = "language" //
)

func PlayerInfoKey(uid uint64) string {
	return "player:" + util.Uint64ToString(uid)
}
func GetPlayerBaseInfo(uid uint64, filed string) string {
	return redis.HGet(PlayerInfoKey(uid), filed)
}
func GetPlayerBaseInfos(uid uint64, filed []string) map[string]string {
	return redis.HMGetMap("player:"+util.Uint64ToString(uid), filed)
}
func SetPlayerBaseInfo(uid uint64, field interface{}, value interface{}) {
	redis.HSet("player:"+util.Uint64ToString(uid), field, value)
}

func StorePlayerInfo(uid uint64, info *protoMsg.UserInfo) {
	if info == nil {
		return
	}

	if info.GetNickName() != "" {
		redis.HSet(PlayerInfoKey(uid), NickNameField, info.GetNickName())
	}

	if info.GetAvatarUrl() != "" {
		redis.HSet(PlayerInfoKey(uid), PictureField, info.GetAvatarUrl())
	}

	if info.GetCountry() != "" {
		redis.HSet(PlayerInfoKey(uid), CountryField, info.GetCountry())
	}

	if info.GetProvince() != "" {
		redis.HSet(PlayerInfoKey(uid), ProvinceField, info.GetProvince())
	}

	if info.GetCity() != "" {
		redis.HSet(PlayerInfoKey(uid), CityField, info.GetCity())
	}

	if info.GetLanguage() != "" {
		redis.HSet(PlayerInfoKey(uid), LanguageField, info.GetLanguage())
	}

	if info.GetGender() != "" {
		redis.HSet(PlayerInfoKey(uid), GenderField, info.GetGender())
	}
}
