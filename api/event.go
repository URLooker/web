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

func (this *Web) SavePortEvent(event *model.PortEvent, reply *string) error {
	err := event.Insert()
	if err != nil {
		*reply = err.Error()
	}
	return nil
}