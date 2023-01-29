package infra

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/111161226/TOKO-ENCOUNT/db/model"
	"github.com/111161226/TOKO-ENCOUNT/db/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userInfra struct {
	db *sqlx.DB
}

func NewUserInfra(db *sqlx.DB) repository.UserRepository {
	return &userInfra{db: db}
}

//func for hashing password
func hash(pw string) string {
	const salt = "SakataKintoki#"
	h := sha256.New()
	h.Write([]byte(salt))
	h.Write([]byte(pw))
	return hex.EncodeToString(h.Sum(nil))
}

//delete user by logic
func (ui *userInfra) DeleteUser(userId string) error {
	//get username by userId
	var username string
	err := ui.db.Get(
		&username, "SELECT `user_name` FROM `users` WHERE `user_id` = ?",
		userId,
	)
	if err != nil {
		return err
	}
	//delete user by username
	_, err = ui.db.Exec(
		"UPDATE `user_deletes` SET `flag` = ? WHERE `user_name` = ?",
		1,
		username,
	)
	if err != nil {
		return err
	} 
	return err
}

//check user data was used
func (ui *userInfra) CheckUsedUser(userName string, password string) (*model.UserWithoutPass, error) {
	//get user data
	var user model.UserWithoutPass
	err := ui.db.Get(
		&user, "SELECT `user_id`, `user_name`, `prefect`, `gender` FROM `users` WHERE `user_name` = ? AND `password` = ?",
		userName,
		hash(password),
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//restore user
func (ui *userInfra) RestoreUser(userId string) error {
	//update DB
	_, err := ui.db.Exec(
		"UPDATE `user_deletes` SET `flag` = ? WHERE `user_id` = ?",
		1,
		userId,
	)
	if err != nil {
		return  err
	}
	return nil
}

//create user
func (ui *userInfra) CreateUser(user *model.User) (*model.UserWithoutPass, error) {
	//set uuid
	uu, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	userId := uu.String()

	//insert data into DB
	_, err = ui.db.Exec(
		"INSERT INTO `users`(`user_id`, `user_name`, `password`, `prefect`, `gender`) VALUES (?, ?, ?, ? ,?)",
		userId,
		user.UserName,
		hash(user.Password),
		user.Prefect,
		user.Gender,
	)
	if err != nil {
		return nil, err
	}
	_, err = ui.db.Exec(
		"INSERT INTO `user_deletes` (`user_name`) VALUES (?)",
		user.UserName,
	)
	if err != nil {
		return nil, err
	}

	return &model.UserWithoutPass{
		UserId:   userId,
		UserName: user.UserName,
		Prefect:  user.Prefect,
		Gender:   user.Gender,
	}, nil
}

func (ui *userInfra) GetUser(userId string) (*model.UserWithoutPass, error) {
	//get user data
	var user model.UserWithoutPass
	err := ui.db.Get(
		&user, "SELECT `user_id`, `user_name`, `prefect`, `gender` FROM `users` WHERE `user_id` = ?",
		userId,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ui *userInfra) EditUser(userId string, user *model.UserUpdate) (*model.UserWithoutPass, error) {
	//get old password
	var oldpassword string
	err := ui.db.Get(
		&oldpassword, "SELECT `password` FROM `users` WHERE `user_id` = ?",
		userId,
	)
	if err != nil {
		return nil, err
	}

	//check password is right
	if oldpassword != hash(user.Password) {
		return nil, nil //in case password is incorrect , return nil
	}

	//update DB
	_, err = ui.db.Exec(
		"UPDATE `users` SET `user_name` = ?, `password` = ?, `prefect` = ?, `gender` = ? WHERE `user_id` = ?",
		user.UserName,
		hash(user.NewPassword),
		user.NewPrefect,
		user.NewGender,
		userId,
	)
	if err != nil {
		return nil, err
	}

	return &model.UserWithoutPass{
		UserId:   userId,
		UserName: user.UserName,
		Prefect:  user.NewPrefect,
		Gender:   user.NewGender,
	}, nil
}

func (ui *userInfra) CheckRightUser(user *model.UserSimple) (*model.UserWithoutPass, error) {
	//get password
	var password string
	err := ui.db.Get(
		&password, "SELECT `password` FROM `users` WHERE `user_name` = ?",
		user.UserName,
	)
	if err != nil {
		return nil, err
	}

	//check password is right
	if password != hash(user.Password) {
		return nil, fmt.Errorf("err : %s", "Incorrect password")
	}

	//get user data
	var userwithoutpass model.UserWithoutPass
	err = ui.db.Get(
		&userwithoutpass, "SELECT `user_id`, `user_name`, `prefect`, `gender` FROM `users` WHERE `user_name` = ?",
		user.UserName,
	)
	if err != nil {
		return nil, err
	}

	return &userwithoutpass, nil
}

func (ui *userInfra) CheckUsedUserName(userName string) (*model.UserWithoutPass, error) {
	//get userName
	var user model.UserWithoutPass
	err := ui.db.Get(
		&user, "SELECT `user_id`, `user_name`, `prefect`, `gender` FROM `users` WHERE `user_name` = ?",
		userName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//in case no one uses the username, that is not error 
			return nil, nil
		}

		return nil, err
	}

	//in case the username is duplicated, user is not nil
	return &user, nil
}

func (ui *userInfra) GetUserList(limit int, offset int, name string, gender string, prefect string, user_id string) (*model.UserList, error) {
	//create query
	query := ""
	query1 := "SELECT `user_id`, `user_name`, `prefect`, `gender` FROM `users` WHERE `user_id` != ? "
	query2 := "SELECT COUNT(*) FROM `users` WHERE `user_id` != ? "
	bind := []interface{}{
		user_id,
	}
	//check condition is added 
	if name != "" {
		query += "AND `user_name` LIKE ? "
		bind = append(bind, "%"+name+"%")
	}

	if gender != "" {
		query += "AND `gender` = ? "
		bind = append(bind, gender)
	}

	if prefect != "" {
		query += "AND `prefect` = ? "
		bind = append(bind, prefect)
	}

	//get designated user 
	query1 = query1 + query + "LIMIT ? OFFSET ? "
	bind1 := append(bind, limit, offset)
	users := []*model.UserWithoutPass{}
	err := ui.db.Select(
		&users,
		query1,
		bind1...,
	)
	if err != nil {
		return nil, err
	}

	//get the number of the designated user
	query2 = query2 + query
	var count int
	err = ui.db.Get(
		&count,
		query2,
		bind...,
	)

	return &model.UserList{
		HasNext: count > len(users)+offset,
		Users:   &users,
	}, nil
}
