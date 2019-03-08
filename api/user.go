package api

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"

	"github.com/peng19940915/urlooker/web/model"
)

type UsersResponse struct {
	Message string
	Data    []*model.RelTeamUser
}

func (this *Web) GetUsersByTeam(req string, reply *UsersResponse) error {
	tids := strings.Split(req, ",")
	if len(tids) < 1 || tids[0] == "" {
		reply.Message = "user no exists!"
		return nil
	}
	allUsers := make([]*model.RelTeamUser, 0)
	for _, tid := range tids {
		id, err := strconv.ParseInt(tid, 10, 64)
		if err != nil {
			log.Errorf("tid error,detail: %v", err)
			continue
		}
		users, err := model.MailsOfTeam(id)
		if err != nil {
			reply.Message = err.Error()
		}

		for _, user := range users {
			allUsers = append(allUsers, user)
		}
	}

	reply.Data = allUsers
	return nil
}
