package api

import (
	"github.com/urlooker/web/model"
)

func (this *Web) SaveEvent(event *model.Event, reply *string) error {
	err := event.Insert()
	if err != nil {
		*reply = err.Error()
	}

	return nil
}
