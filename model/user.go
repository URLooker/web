package model

import (
	"fmt"
	"log"
	"time"

	. "github.com/urlooker/web/store"
)

type User struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Cnname   string    `json:"cnname"`
	Password string    `json:"-"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Wechat   string    `json:"wechat"`
	Role     int       `json:"role"`
	Created  time.Time `json:"-" xorm:"<-"`
}

func GetUserById(id int64) (*User, error) {
	var obj User
	has, err := Orm.Cols("id", "name", "cnname", "email", "phone", "wechat").Where("id=?", id).Get(&obj)
	if err != nil || !has {
		return nil, err
	}
	obj.Password = ""
	return &obj, nil
}

// 先根据name获取id，再根据id去查询
func GetUserByName(name string) (*User, error) {
	if name == "" {
		return nil, nil
	}

	var obj User
	has, err := Orm.Cols("id", "name", "cnname", "email", "phone", "wechat").Where("name=?", name).Get(&obj)
	if err != nil || !has {
		return nil, err
	}

	return &obj, nil
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

func UserRegister(name, password string) (int64, error) {
	user, err := GetUserByName(name)
	if err != nil {
		return 0, err
	}
	if user != nil {
		return 0, fmt.Errorf("用户名已经被占用了")
	}

	user = &User{Name: name, Password: password}
	_, err = Orm.Insert(user)
	if err != nil {
		return 0, err
	}

	return user.Id, nil
}

func UserLogin(name, password string) (int64, error) {
	user, err := GetUserByName(name)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, fmt.Errorf("系统中没有该用户%s", name)
	}
	if user.Password != password {
		return 0, fmt.Errorf("密码不正确")
	}

	return user.Id, nil
}

func (this *User) UpdateProfile() error {
	log.Println(this)
	_, err := Orm.Id(this.Id).Update(this)
	if err != nil {
		return err
	}

	return err
}

func (this *User) ChangePasswd(oldPasswd, newPasswd string) error {
	if this.Password != oldPasswd {
		return fmt.Errorf("原密码不正确")
	}

	this.Password = newPasswd
	_, err := Orm.Id(this.Id).Update(this)
	if err != nil {
		return err
	}

	return nil
}
