package model

import (
. "github.com/peng19940915/urlooker/web/store"
)

type PortStrategy struct {
	Id         int64  `json:"id"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Enable     int    `json:"enable"`
	IP         string `json:"ip" xorm:"ip"`
	Keywords   string `json:"keywords"`
	Timeout    int    `json:"timeout"`
	Creator    string `json:"creator"`
	Note       string `json:"note"`
	Tag        string `json:"tag"`
	MaxStep    int    `json:"max_step"`
	Times      int    `json:"times"`
	Teams      string `json:"teams"`
}

func GetAllPortStrategyCount(mine int, query, username string) (int64, error) {
	if mine == 1 {
		if query != "" {
			return Orm.Where("host LIKE ? AND creator = ? ORDER BY id", "%"+query+"%", username).Count(new(PortStrategy))
		} else {
			num, err := Orm.Where("creator = ?", username).Count(new(PortStrategy))
			return num, err
		}
	} else {
		if query != "" {
			return Orm.Where("host LIKE ? ORDER BY id", "%"+query+"%").Count(new(PortStrategy))
		} else {
			num, err := Orm.Count(new(PortStrategy))
			return num, err
		}
	}

}

func GetAllPortStrategy(mine, limit, offset int, query, username string) ([]*PortStrategy, error) {
	items := make([]*PortStrategy, 0)

	var err error
	if mine == 1 {
		if query != "" {
			err = Orm.Where("host LIKE ? AND creator = ? ORDER BY id", "%"+query+"%", username).Limit(limit, offset).Find(&items)
		} else {
			err = Orm.Where("creator = ?", username).Limit(limit, offset).Find(&items)
		}
	} else {
		if query != "" {
			err = Orm.Where("host LIKE ? ORDER BY id", "%"+query+"%").Limit(limit, offset).Find(&items)
		} else {
			err = Orm.Limit(limit, offset).Find(&items)
		}
	}
	return items, err
}

func GetAllPortStrategyByCron() ([]*PortStrategy, error) {
	strategies := make([]*PortStrategy, 0)
	err := Orm.Where("enable = 1").Find(&strategies)

	return strategies, err
}

func GetPortStrategyById(sid int64) (*PortStrategy, error) {
	strategy := new(PortStrategy)
	_, err := Orm.Where("id=?", sid).Get(strategy)
	return strategy, err
}

func (this *PortStrategy) Add() (int64, error) {
	_, err := Orm.Insert(this)
	return this.Id, err
}

func (this *PortStrategy) Update() error {
	_, err := Orm.Where("id=?", this.Id).Cols("times", "max_step", "port", "timeout", "host", "enable", "ip", "keywords", "note", "data", "tag", "teams").Update(this)
	return err
}

func (this *PortStrategy) Delete() error {
	_, err := Orm.Where("id=?", this.Id).Delete(new(PortStrategy))
	return err
}
