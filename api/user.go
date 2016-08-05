package api

import (
	"log"
	"strconv"
	"strings"

	"github.com/urlooker/web/model"
)

type UsersResponse struct {
	Message string
	Data    []*model.User
}

func (this *Web) GetUsersByTeam(req string, reply *UsersResponse) error {
	tids := strings.Split(req, ",")
	if len(tids) < 1 || tids[0] == "" {
		reply.Message = "user no exists!"
		return nil
	}
	allUsers := make([]*model.User, 0)
	for _, tid := range tids {
		id, err := strconv.ParseInt(tid, 10, 64)
		if err != nil {
			log.Println("tid error:", err)
			continue
		}
		users, err := model.UsersInfoOfTeam(id)
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
