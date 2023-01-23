package repository

import "github.com/cs-sysimpl/SakataKintoki/db/model"

type UserRepository interface {
	CreateUser(user *model.User) (*model.UserWithoutPass, error)
	GetUser(userId string) (*model.UserWithoutPass, error)
	EditUser(userId string, user *model.UserUpdate) (*model.UserWithoutPass, error)
	//withdraw user by logic delete
	//DeleteUser(userId string, pass string) error
	//check whether user can login
	CheckRightUser(user *model.UserSimple) (*model.UserWithoutPass, error)
	//check username is duplicated
	CheckUsedUserName(userName string) (*model.UserWithoutPass, error)
	//search users
	GetUserList(limit int, offset int, name string, gender string, prefect string, user_id string) (*model.UserList, error)
}
