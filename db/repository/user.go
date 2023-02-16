package repository

import "github.com/111161226/TOKO-ENCOUNT/db/model"

type UserRepository interface {
	CreateUser(user *model.User) (*model.UserWithoutPass, error)
	GetUserByUserId(userId string) (*model.UserWithoutPass, error)
	GetUserByUserName(userName string) (*model.UserWithoutPass, error)
	EditUser(userId string, user *model.UserUpdate) (*model.UserWithoutPass, error)
	//withdraw user by logic delete
	DeleteUser(userId string) error
	//check user data was used
	CheckUsedUser(userName string, password string) (*model.UserWithoutPass, error)
	//restore user data
	RestoreUser(userName string) error
	//check whether user can login
	CheckRightUser(user *model.UserSimple) (*model.UserWithoutPass, error)
	//check username is duplicated
	CheckUsedUserName(userName string) (*model.UserWithoutPass, error)
	//search users
	GetUserList(limit int, offset int, name string, gender string, prefect string, user_id string) (*model.UserList, error)
	GetRoomUsers(roomId string, userId string) (*model.UserList, error)
}
