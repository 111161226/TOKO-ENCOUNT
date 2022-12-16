package repository

import "github.com/cs-sysimpl/SakataKintoki/db/model"

type UserRepository interface {
	CreateUser(user *model.UserCreate) (*model.UserWithoutPass, error)
	GetUser(userId string) (*model.UserWithoutPass, error)
	EditUser(userId string, user *model.UserUpdate) (*model.UserWithoutPass, error)
	//withdraw user by logic delete
	DeleteUser(userId string, pass string) error
	//check whether user can login
	CheckRightUser(user *model.UserLogin) (*model.UserWithoutPass, error)
	//check username is duplicated
	CheckUsedUserName(userName string) (*model.UserWithoutPass, error)
}