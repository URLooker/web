package cron

import (
	"log"
	"time"

	"github.com/urlooker/web/model"
)

func DeleteOld() {
	d := 12
	t1 := time.NewTicker(time.Duration(d) * time.Second)
	for {
		<-t1.C
		err := model.RelSidIpRepo.DeleteOld(int64(d))
		if err != nil {
			log.Println("delete error:", err)
		}

		err = model.ItemStatusRepo.DeleteOld(int64(d))
		if err != nil {
			log.Println("delete error:", err)
		}
	}
}
