package infra

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/cs-sysimpl/SakataKintoki/db/model"
	"github.com/cs-sysimpl/SakataKintoki/db/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userInfra struct {
	db *sqlx.DB
}

func NewUserInfra(db *sqlx.DB) repository.UserRepository {
	return &userInfra{db: db}
}

//ハッシュ化関数
func hash(pw string) string {
	const salt = "SakataKintoki#"
	h := sha256.New()
	h.Write([]byte(salt))
	h.Write([]byte(pw))
	return hex.EncodeToString(h.Sum(nil))
}

func (ui *userInfra) CreateUser(user *model.User) (*model.UserWithoutPass, error) {
	//UUID設定
	uu, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	userId := uu.String()

	//DB挿入
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

	return &model.UserWithoutPass{
		UserId:   userId,
		UserName: user.UserName,
		Prefect:  user.Prefect,
		Gender:   user.Gender,
	}, nil
}

func (ui *userInfra) GetUser(userId string) (*model.UserWithoutPass, error) {
	//ユーザ取得
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
	//古いパスワード取得
	var oldpassword string
	err := ui.db.Get(
		&oldpassword, "SELECT `password` FROM `users` WHERE `user_id` = ?",
		userId,
	)
	if err != nil {
		return nil, err
	}

	//パスワード照合
	if oldpassword != hash(user.Password) {
		return nil, nil //間違っている場合は返り値nil
	}

	//DB更新
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
	//パスワード取得
	var password string
	err := ui.db.Get(
		&password, "SELECT `password` FROM `users` WHERE `user_name` = ?",
		user.UserName,
	)
	if err != nil {
		return nil, err
	}

	//パスワード照合
	if password != hash(user.Password) {
		return nil, fmt.Errorf("err : %s", "Incorrect password")
	}

	//User取得
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
	// userName取得
	var user model.UserWithoutPass
	err := ui.db.Get(
		&user, "SELECT `user_id`, `user_name`, `prefect`, `gender` FROM `users` WHERE `user_name` = ?",
		userName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 名前の重複したユーザーが存在しないならエラーではない
			return nil, nil
		}

		return nil, err
	}

	// 重複している場合はuser!=nil
	return &user, nil
}

func (ui *userInfra) GetUserList(limit int, offset int, name string, gender string, prefect string) (*model.UserList, error) {
	//対象となるユーザを取得
	var users []*model.UserWithoutPass
	err := ui.db.Select(
		&users, "SELECT `user_id`, `user_name`, `prefect`, `gender` FROM `users` WHERE CASE WHEN ? = ? THEN ? ELSE `user_name` END = ? AND CASE WHEN ? = ? THEN ? ELSE `gender` END = ? AND CASE WHEN ? = ? THEN ? ELSE `prefect` END = ? LIMIT ? OFFSET ?",
		name,
		"",
		name,
		name,
		gender,
		"",
		gender,
		gender,
		prefect,
		"",
		prefect,
		prefect,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	//対象となるユーザの全数を取得
	var count int
	err = ui.db.Get(
		&count, "SELECT COUNT(*) FROM `users` WHERE CASE WHEN ? = ? THEN ? ELSE `user_name` END = ? AND CASE WHEN ? = ? THEN ? ELSE `gender` END = ? AND CASE WHEN ? = ? THEN ? ELSE `prefect` END = ?",
		name,
		"",
		name,
		name,
		gender,
		"",
		gender,
		gender,
		prefect,
		"",
		prefect,
		prefect,
	)
	if err != nil {
		return nil, err
	}

	return &model.UserList{
		HasNext: count > len(users) + offset,
		Users: &users,
	}, nil
}