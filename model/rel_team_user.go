package model

import (
	"github.com/go-xorm/xorm"
	. "github.com/peng19940915/urlooker/web/store"
)

type RelTeamUser struct {
	Id    int64  `json:"id"`
	Tid   int64  `json:"tid"`
	Email string `json:"email"`
}

func Teams(query string,limit, offset int) ([]*Team, error) {
	teams := make([]*Team, 0)
	var err error
	if query != "" {
		tn := "%" + query + "%"
		err = Orm.Sql("SELECT * FROM team WHERE name LIKE ? ORDER BY name LIMIT ?,?", tn, offset, limit).Find(&teams)
	} else {
		err = Orm.Sql("SELECT * FROM team ORDER BY name LIMIT ?,?", offset, limit).Find(&teams)
	}

	return teams, err
}

func TeamCountOfUser(query string) (int64, error) {
	if query != "" {
		tn := "%" + query + "%"
		return Orm.Where("email LIKE ? ", tn).Count(new(Team))
	} else {
		return Orm.Count(new(Team))
	}
}

func MailsOfTeam(tid int64) ([]*RelTeamUser, error) {
	var err error
	users := make([]*RelTeamUser, 0)
	if tid == 0 {
		err = Orm.Cols("email").Sql("SELECT * FROM rel_team_user WHERE tid=?", tid).Find(&users)
	}else{
		err = Orm.Cols("email").Sql("SELECT * FROM rel_team_user WHERE tid=?", tid).Find(&users)
	}

	if err != nil {
		return users, err
	}
	return users, nil
}



func IsCreatorOfTeam(uid, tid int64) (bool, error) {
	team, err := GetTeamById(tid)
	if err != nil {
		return false, err
	}

	if team.Creator == uid {
		return true, nil
	}
	return false, nil
}

func updateUsersOfTeam(tid int64, emails []string, session *xorm.Session) error {
	err := removeAllUsersFromTeam(tid, session)
	if err != nil {
		return err
	}

	err = addUsersIntoTeam(tid, emails, session)
	if err != nil {
		return err
	}

	return err
}

func addUsersIntoTeam(tid int64, emails []string, session *xorm.Session) error {

	if len(emails) == 0 {
		relTeamUser := &RelTeamUser{Tid: tid}
		_, err := session.Insert(relTeamUser)
		if err != nil {
			return err
		}

	}
	for _, email := range emails {
		relTeamUser := &RelTeamUser{Tid: tid, Email: email}
		_, err := session.Insert(relTeamUser)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeAllUsersFromTeam(tid int64, session *xorm.Session) error {
	_, err := session.Where("tid=?", tid).Delete(new(RelTeamUser))
	return err
}
