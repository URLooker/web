package api

import (
	"github.com/peng19940915/urlooker/web/model"
)

func (this *Web) SaveEvent(event *model.Event, reply *string) error {
	err := event.Insert()
	if err != nil {
		*reply = err.Error()
	}
	return nil
}

type EventCacheReply struct {
	Message string
	Data [] *model.Event
}
func (this *Web) GetEvent(arg string, reply *EventCacheReply) error {
	items, err := model.GetAlarmEvent()
	if err != nil {
		reply.Message = "Get history Alarm from database failed, detail: "+ err.Error()
	}
	reply.Data=items
	return nil

}


func (this *Web) SavePortEvent(event *model.PortEvent, reply *string) error {
	err := event.Insert()
	if err != nil {
		*reply = err.Error()
	}
	return nil
}


type PortEventCacheReply struct {
	Message string
	Data [] *model.PortEvent
}
func (this *Web) GetPortEvent(arg string, reply *PortEventCacheReply) error {
	items, err := model.GetPortAlarmEvent()
	if err != nil {
		reply.Message = "Get history Alarm from database failed, detail: "+ err.Error()
	}
	reply.Data=items
	return nil

}

