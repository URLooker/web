package model

import (
	. "github.com/peng19940915/urlooker/web/store"
	"fmt"
	"time"
	"github.com/peng19940915/urlooker/web/g"
	"crypto/md5"
	"io"
)

type PortStatus struct {
	Id       int64  `json:"id"`
	Sid      int64  `json:"sid"`
	RespTime int    `json:"resp_time"`
	PushTime int64  `json:"push_time"`
	Result   int64    `json:"result"`
}

var PortStatusRepo *PortStatus

func (this *PortStatus) Save() error {
	sql := fmt.Sprintf("insert into port_status00 (id, sid, resp_time, push_time, result) value(?,?,?,?,?)")
	_, err := Orm.Exec(sql, this.Id, this.Sid, this.RespTime, this.PushTime, this.Result)
	return err
}

func (this *PortStatus) GetBySid(sid int64) ([] *PortStatus, error) {
	portStatusArr := make([] *PortStatus, 0)
	ts := time.Now().Unix() - int64(g.Config.Past * 60)
	sql := fmt.Sprintf("select * from port_status00 where sid=? and push_time > ?")
	err := Orm.Sql(sql, sid, ts).Find(&portStatusArr)
	return portStatusArr, err
}

func (this *PortStatus) PK() string {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%s", this.Sid))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (this *PortStatus) DeleteOld(d int64) error {
	ts := time.Now().Unix() - 12*60*60
	sql := fmt.Sprintf("delete from port_status00 where push time < ?")
	_, err := Orm.Exec(sql, ts)
	return err
}

