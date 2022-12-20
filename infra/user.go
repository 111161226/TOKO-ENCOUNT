package infra

import (
	"github.com/cs-sysimpl/SakataKintoki/db/model"
	"github.com/cs-sysimpl/SakataKintoki/db/repository"
	"github.com/jmoiron/sqlx"
)

type userInfra struct {
	db *sqlx.DB
}

func NewUserInfra(db *sqlx.DB) repository.UserRepository {
	return &userInfra{db: db}
}

func (ui *userInfra) CreateUser(user *model.User) (*model.UserWithoutPass, error) {
	return nil, nil
}

func (ui *userInfra) GetUser(userId string) (*model.UserWithoutPass, error) {
	return nil, nil
}

func (ui *userInfra) EditUser(userId string, user *model.UserUpdate) (*model.UserWithoutPass, error) {
	return nil, nil
}

func (ui *userInfra) CheckRightUser(user *model.UserSimple) (*model.UserWithoutPass, error) {
	return nil, nil
}

func (ui *userInfra) CheckUsedUserName(userName string) (*model.UserWithoutPass, error) {
	return nil, nil
}
