package user

import (
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

func CheckUser(env *beldilib.Env, req Request) Result {
	var user data.User
	item := beldilib.Read(env, data.Tuser(), req.Username)
	err := mapstructure.Decode(item, &user)
	beldilib.CHECK(err)
	return Result{Correct: req.Password == user.Password}
}
