package repositoty

import (
	"backend_myblog/model"
	"backend_myblog/model/req"
	"context"
)

type UserRepo interface {
	CheckLogin(context context.Context, loginReq  req.ReqSignIn) (model.User , error)
	SaveUser(context context.Context, user model.User) (model.User , error)
	SelectUserByID(context context.Context, userId string) (model.User , error)
	UpdateUser(context context.Context, user model.User) (model.User, error)
}