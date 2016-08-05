package model

import (
	"time"

	"github.com/urlooker/web/g"
	. "github.com/urlooker/web/store"
)

type RelSidIp struct {
	Id  int64  `json:"id"`
	Sid int64  `json:"sid"`
	Ip  string `json:"ip"`
	Ts  int64  `json:"ts"`
}

var RelSidIpRepo *RelSidIp

func (this *RelSidIp) Save() error {
	has, err := Orm.Where("sid = ? and ip = ?", this.Sid, this.Ip).Get(new(RelSidIp))
	if err != nil {
		return err
	}
	if !has {
		_, err = Orm.Insert(this)
	} else {
		this.Ts = time.Now().Unix()
		_, err = Orm.Cols("ts").Update(this)
	}
	return err
}

func (this *RelSidIp) GetBySid(sid int64) ([]*RelSidIp, error) {
	var relSidIps []*RelSidIp
	ts := time.Now().Unix() - int64(g.Config.Past*60)
	err := Orm.Where("sid = ? and ts > ?", sid, ts).Find(&relSidIps)
	return relSidIps, err
}

func (this *RelSidIp) DeleteOld(d int64) error {
	ts := time.Now().Unix() - d*60*60
	_, err := Orm.Where("ts < ?", ts).Delete(new(RelSidIp))
	return err
}
