package cron

import (
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/peng19940915/urlooker/web/model"
)

func DeleteOld() {
	d := 12
	t1 := time.NewTicker(time.Duration(d) * time.Hour)
	for {
		<-t1.C
		err := model.RelSidIpRepo.DeleteOld(int64(d))
		if err != nil {
			log.Errorf("delete error: %v", err)
		}
		// 清理过期port 历史数据
		err = model.PortStatusRepo.DeleteOld(int64(d))
		if err != nil {
			log.Errorf("delete error: %v", err)
		}
		// 清理过期url 历史数据（12h）
		err = model.ItemStatusRepo.DeleteOld(int64(d))
		if err != nil {
			log.Errorf("delete error: %v", err)
		}
		// 清理过期port event (30d)
		err = model.DeleteOldPortEvent()
		if err != nil {
			log.Errorf("delete error: %v", err)
		}

		err = model.DeleteOldEvent()
		if err != nil {
			log.Errorf("delete error: %v", err)
		}
	}
}
