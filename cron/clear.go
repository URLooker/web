package cron

import (
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/peng19940915/urlooker/web/model"
)

func DeleteOld() {
	d := 12
	t1 := time.NewTicker(time.Duration(d) * time.Second)
	for {
		<-t1.C
		err := model.RelSidIpRepo.DeleteOld(int64(d))
		if err != nil {
			log.Errorf("delete error: %v", err)
		}

		err = model.ItemStatusRepo.DeleteOld(int64(d))
		if err != nil {
			log.Errorf("delete error: %v", err)
		}
	}
}
