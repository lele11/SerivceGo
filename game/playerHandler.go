package game

import (
	"game/db"
	"game/protoMsg"
)

func (pm *PlayerManager) AttachPlayerHandler() {
	pm.AttachHandler(protoMsg.C_CMD_C_ACCOUNTINFO, pm.HandlerAccountInfo)

}
func (pm *PlayerManager) HandlerAccountInfo(p *Player, data []byte) {
	ret := &protoMsg.S_AccountInfo{}
	info := db.GetPlayerBaseInfos(p.uid, []string{
		db.ChannelField,
		db.NickNameField,
		db.PictureField,
		db.GenderField,
		db.CountryField,
		db.ProvinceField,
		db.CityField,
		db.LanguageField,
	})

	ret.Uid = p.uid
	ret.Picture = info[db.PictureField]
	ret.Gender = info[db.GenderField]
	ret.NickName = info[db.NickNameField]
	ret.Country = info[db.CountryField]
	ret.Province = info[db.ProvinceField]
	ret.City = info[db.CityField]
	ret.Language = info[db.LanguageField]

	p.SendToClient(protoMsg.S_CMD_S_ACCOUNTINFO, ret)
}
