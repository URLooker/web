package model

import (
	"fmt"
	"time"

	. "github.com/peng19940915/urlooker/web/store"
)

type Event struct {
	Id           int64  `json:"id"`
	EventId      string `json:"event_id"`
	Status       string `json:"status"`
	Url          string `json:"url"`
	Ip           string `json:"ip"`
	EventTime    int64  `json:"event_time"`
	StrategyId   int64  `json:"strategy_id"`
	RespTime     int    `json:"resp_time"`
	RespCode     string `json:"resp_code"`
	Result       int64  `json:"result"`
	CurrentStep  int    `json:"current_step"`
	MaxStep      int    `json:"max_step"`
}

var EventRepo *Event

func (this *Event) Insert() error {
	if this.Status == "OK" {
		_, err := Orm.Exec(`update event set status=? where event_id= ?`,"OK", this.EventId)
		return err
	}
	_, err := Orm.Insert(this)
	return err
}

func (this *Event) Update()error {
	_,err := Orm.Update(this)
	return err
}

func (this *Event) GetByStrategyId(strategyId int64, before int) ([]*Event, error) {
	events := make([]*Event, 0)
	ts := time.Now().Unix() - int64(before)
	err := Orm.Where("strategy_id = ? and event_time > ?", strategyId, ts).Desc("event_time").Find(&events)
	return events, err
}

func (this *Event) String() string {
	return fmt.Sprintf(
		"<Id:%s, EventId:,%s Ip:%s, Url:%s, EventTime:%v, StrategyId:%d, RespTime:%s, RespCode:%s, Status:%s (%d/%d)>",
		this.Id,
		this.EventId,
		this.Ip,
		this.Url,
		this.EventTime,
		this.StrategyId,
		this.RespTime,
		this.RespCode,
		this.Status,
		this.CurrentStep,
		this.MaxStep,
	)
}


func GetAllEventCounter(query string)(int64, error){

	if query != "" {
		return Orm.Where("url LIKE ? ", "%"+query+"%").Count(new(Event))
	}else {
		num, err := Orm.Desc("status").Count(new(Event))
		return num, err
	}
}

func GetAllEvent(limit, offset int, query string)([]*Event, error) {
	items := make([]*Event, 0)
	var err error
	if query != "" {
		err = Orm.Where("url like ? ", "%"+query+"%").Limit(limit, offset).Desc("status").Find(&items)

	}else {
		err = Orm.Desc("status").Limit(limit, offset).Find(&items)
	}
	return items, err
}

