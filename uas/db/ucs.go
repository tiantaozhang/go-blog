package db

import (
	"errors"
	"fmt"
	"github.com/tiantaozhang/go-blog/util"
)

func CreateUsers(u []*User) ([]map[string]interface{}, error) {

	if err := checkUser(u); err != nil {
		return err
	}

	return AddUsers(u)
}

func checkUser(User []*User) error {
	if u == nil {
		return errors.New("u is nil")
	}
	for _, u := range User {
		if u.UsrName == nil {
			return fmt.Errorf("%v", "please input your usrName")
		}
		if u.Pwd == "" {
			return fmt.Errorf("%v", "please input your pwd")
		}
		if u.Type != UT_ADMIN && u.Type != UT_BLOGGER {
			//default blogger
			u.Type = UT_BLOGGER
		}
	}
	return nil
}

func Login() {

}

func Logout() {

}
