package helper

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"

	"github.com/peng19940915/urlooker/web/model"
)

func Times1000(num int64) int64 {
	return (num + 3600*8) * 1000
}

func UsersOfTeam(tid int64) []*model.User {
	users, _ := model.UsersOfTeam(tid)
	return users
}

func TeamsOfStrategy(ids string) []*model.Team {
	teams, err := model.GetTeamsByIds(ids)
	if err != nil {
		log.Errorf("get teams err, detail:%v", err.Error())
	}
	return teams
}

func HumenTime(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func GetFirst(items []*model.ItemStatus) string {
	if len(items) == 0 {
		return " , "
	}
	item := items[0]
	str := HumenTime(item.PushTime) + "," + strconv.Itoa(item.RespTime)
	return str
}
