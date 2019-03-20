package model

import (
	. "github.com/peng19940915/urlooker/web/store"
	"time"
	"fmt"
)

type PortEvent struct {
	Id          int64  `json:"id"`
	EventId     string `json:"event_id"`
	Status      string `json:"status"`
	Host        string `json:"host"`
	Ip          string `json:"ip"`
	Port        string `json:"port"`
	EventTime   int64  `json:"event_time"`
	StrategyId  int64  `json:"strategy_id"`
	RespTime    int    `json:"resp_time"`
	Result      int64  `json:"result"`
	CurrentStep int    `json:"current_step"`
	MaxStep     int    `json:"max_step"`
}

var PortEventRepo * PortEvent

func (this *PortEvent) Insert() error {
	if this.Status == "OK" {
		_, err := Orm.Exec(`update port_event set status=? where event_id= ?`,"OK", this.EventId)

		return err
	}
	_, err := Orm.Insert(this)
	return err
}

func (this *PortEvent) Update()error {
	_,err := Orm.Update(this)
	return err
}

func (this *PortEvent) GetByPortStrategyId(strategyId int64, before int)([]*PortEvent, error) {
	events := make([] *PortEvent, 0)
	ts := time.Now().Unix() - int64(before)
	err := Orm.Where("strategy_id = ? and event_time > ?", strategyId, ts).Desc("event_time").Find(&events)
	return events, err
}

func (this *PortEvent) String() string {
	return fmt.Sprintf(
		"<Id:%s, EventId:,%s Ip:%s, Port:%s, Host:%s, EventTime:%v, StrategyId:%d, RespTime:%s, RespCode:%s, Status:%s, (%d/%d)>",
		this.Id,
		this.EventId,
		this.Ip,
		this.Port,
		this.Host,
		this.EventTime,
		this.StrategyId,
		this.RespTime,
		this.Status,
		this.CurrentStep,
		this.MaxStep,
	)
}

func GetAllPortEventCounter(query string)(int64, error){
	if query != "" {
		return Orm.Where("host LIKE ?", "%"+query+"%").Desc("event_time").Count(new(Event))
	}else {
		num, err := Orm.Desc("status").Count(new(Event))
		return num, err
	}
}

func GetAllPortEvent(limit, offset int, query string)([]*PortEvent, error) {
	items := make([]*PortEvent, 0)
	var err error
	if query != "" {
		err = Orm.Where("host LIKE ?", "%"+query+"%").Limit(limit).Find(&items)
		fmt.Println(query)
	}else {
		err = Orm.Desc("status").Limit(limit, offset).Find(&items)
	}
	return items, err
}

func DeleteOldPortEvent() error{
	ts := time.Now().Unix() - 2592000
	sql := fmt.Sprintf("delete from port_event where event_time < ?")
	_, err := Orm.Exec(sql, ts)
	return err
}