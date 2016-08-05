package model

import (
	"log"

	. "github.com/urlooker/web/store"
)

type Strategy struct {
	Id         int64  `json:"id"`
	Url        string `json:"url"`
	Keywords   string `json:"keywords"`
	Timeout    int    `json:"timeout"`
	Creator    string `json:"creator"`
	ExpectCode string `json:"expect_code"`
	Note       string `json:"note"`
	Data       string `json:"data"`
	Tag        string `json:"tag"`
	MaxStep    int    `json:"max_step"`
	Times      int    `json:"times"`
	Teams      string `json:"teams"`
}

func GetAllStrategyCount(mine int, query, username string) (int64, error) {
	if mine == 1 {
		if query != "" {
			return Orm.Where("url LIKE ? AND creator = ? ORDER BY id", "%"+query+"%", username).Count(new(Strategy))
		} else {
			num, err := Orm.Where("creator = ?", username).Count(new(Strategy))
			log.Println(err)
			return num, err
		}
	} else {
		if query != "" {
			return Orm.Where("url LIKE ? ORDER BY id", "%"+query+"%").Count(new(Strategy))
		} else {
			log.Println("orm:", Orm)
			num, err := Orm.Count(new(Strategy))
			log.Println(err)
			return num, err
		}
	}

}

func GetAllStrategy(mine, limit, offset int, query, username string) ([]*Strategy, error) {
	items := make([]*Strategy, 0)

	var err error
	if mine == 1 {
		if query != "" {
			err = Orm.Where("url LIKE ? AND creator = ? ORDER BY id", "%"+query+"%", username).Limit(limit, offset).Find(&items)
		} else {
			err = Orm.Where("creator = ?", username).Limit(limit, offset).Find(&items)
		}
	} else {
		if query != "" {
			err = Orm.Where("url LIKE ? ORDER BY id", "%"+query+"%").Limit(limit, offset).Find(&items)
		} else {
			err = Orm.Limit(limit, offset).Find(&items)
		}
	}
	return items, err
}

func GetAllStrategyByCron() ([]*Strategy, error) {
	strategies := make([]*Strategy, 0)
	err := Orm.Find(&strategies)
	return strategies, err
}

func GetStrategyById(sid int64) (*Strategy, error) {
	strategy := new(Strategy)
	_, err := Orm.Where("id=?", sid).Get(strategy)

	return strategy, err
}

func (this *Strategy) Add() (int64, error) {
	_, err := Orm.Insert(this)
	return this.Id, err
}

func (this *Strategy) Update() error {
	_, err := Orm.Where("id=?", this.Id).Cols("times", "max_step", "expect_code", "timeout", "url", "keywords", "note", "data", "tag", "teams").Update(this)
	return err
}

func (this *Strategy) Delete() error {
	_, err := Orm.Where("id=?", this.Id).Delete(new(Strategy))
	return err
}
