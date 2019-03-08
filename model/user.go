package model

import (
	"time"

	. "github.com/peng19940915/urlooker/web/store"
)

type User struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Cnname   string    `json:"cnname"`
	Created  time.Time `json:"-" xorm:"<-"`
}

func GetUserById(id int64) (*User, error) {
	var obj User
	has, err := Orm.Cols("id", "name", "cnname").Where("id=?", id).Get(&obj)
	if err != nil || !has {
		return nil, err
	}
	return &obj, nil
}


// 先根据name获取id，再根据id去查询
func GetUserByName(name string) (*User, error) {
	if name == "" {
		return nil, nil
	}

	var obj User
	has, err := Orm.Cols("id", "name").Where("name=?", name).Get(&obj)
	if err != nil || !has {
		return nil, err
	}

	return &obj, nil
}

func NewUser(name string,cnName string) (id int64, err error){
	user, err := GetUserByName(name)
	//如果该账户以前登陆过
	if user != nil {
		return user.Id,nil
	}else {
		var newUser = &User{
			Name: name,
			Cnname:cnName,
		}
		_, err = Orm.Insert(newUser)
		if err != nil {
			return 0, err
		}
		return newUser.Id, nil
	}

}

func QueryUsers(query string, limit int) ([]*User, error) {
	users := make([]*User, 0)
	if query == "" {
		return users, nil
	}

	err := Orm.Cols("id", "name", "cnname", "email", "phone", "wechat").Where("name LIKE ?", "%"+query+"%").Limit(limit).Find(&users)
	if err != nil {
		return users, err
	}

	return users, nil
}
