package infra

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

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

// func for hashing password
func hash(pw string) string {
	const salt = "SakataKintoki#"
	h := sha256.New()
	h.Write([]byte(salt))
	h.Write([]byte(pw))
	return hex.EncodeToString(h.Sum(nil))
}

// delete user by logic
func (ui *userInfra) DeleteUser(userId string) error {
	var username string
	err := ui.db.Get(
		&username, "SELECT user_name FROM users WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return err
	}
	_, err = ui.db.Exec(
		"UPDATE user_deletes SET flag = $1 WHERE user_name = $2",
		1,
		username,
	)
	return err
}

// check user data was used
func (ui *userInfra) CheckUsedUser(userName string, password string) (*model.UserWithoutPass, error) {
	var user model.UserWithoutPass
	err := ui.db.Get(
		&user, "SELECT users.user_id, users.user_name, users.prefect, users.gender FROM users INNER JOIN user_deletes ON users.user_name = user_deletes.user_name WHERE users.user_name = $1 AND users.password = $2 AND user_deletes.flag = $3",
		userName,
		hash(password),
		1,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// restore user
func (ui *userInfra) RestoreUser(userName string) error {
	_, err := ui.db.Exec(
		"UPDATE user_deletes SET flag = $1 WHERE user_name = $2",
		0,
		userName,
	)
	return err
}

// create user
func (ui *userInfra) CreateUser(user *model.User) (*model.UserWithoutPass, error) {
	uu, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	userId := uu.String()

	_, err = ui.db.Exec(
		"INSERT INTO users(user_id, user_name, password, prefect, gender) VALUES ($1, $2, $3, $4, $5)",
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
		"INSERT INTO user_deletes (user_name) VALUES ($1)",
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

func (ui *userInfra) GetUserByUserId(userId string) (*model.UserWithoutPass, error) {
	var user model.UserWithoutPass
	err := ui.db.Get(
		&user, "SELECT user_id, user_name, prefect, gender FROM users WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ui *userInfra) GetUserByUserName(userName string) (*model.UserWithoutPass, error) {
	var user model.UserWithoutPass
	err := ui.db.Get(
		&user, "SELECT user_id, user_name, prefect, gender FROM users WHERE user_name = $1",
		userName,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ui *userInfra) EditUser(userId string, user *model.UserUpdate) (*model.UserWithoutPass, error) {
	var oldpassword string
	err := ui.db.Get(
		&oldpassword, "SELECT password FROM users WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return nil, err
	}

	if oldpassword != hash(user.Password) {
		return nil, nil
	}

	_, err = ui.db.Exec(
		"UPDATE users SET user_name = $1, password = $2, prefect = $3, gender = $4 WHERE user_id = $5",
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
	var password string
	err := ui.db.Get(
		&password, "SELECT password FROM users WHERE user_name = $1",
		user.UserName,
	)
	if err != nil {
		return nil, err
	}

	if password != hash(user.Password) {
		return nil, fmt.Errorf("err : %s", "Incorrect password")
	}

	var userwithoutpass model.UserWithoutPass
	err = ui.db.Get(
		&userwithoutpass, "SELECT user_id, user_name, prefect, gender FROM users WHERE user_name = $1",
		user.UserName,
	)
	if err != nil {
		return nil, err
	}
	return &userwithoutpass, nil
}

func (ui *userInfra) CheckUsedUserName(userName string) (*model.UserWithoutPass, error) {
	var user model.UserWithoutPass
	err := ui.db.Get(
		&user, "SELECT user_id, user_name, prefect, gender FROM users WHERE user_name = $1",
		userName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// get user by the term
func (ui *userInfra) GetUserList(limit int, offset int, name string, gender string, prefect string, user_id string) (*model.UserList, error) {
	bind := []interface{}{user_id}
	pCount := 1 // プレースホルダのカウンター

	queryCond := ""
	if name != "" {
		pCount++
		queryCond += "AND user_name LIKE $" + strconv.Itoa(pCount) + " "
		bind = append(bind, "%"+name+"%")
	}
	if gender != "" {
		pCount++
		queryCond += "AND gender = $" + strconv.Itoa(pCount) + " "
		bind = append(bind, gender)
	}
	if prefect != "" {
		pCount++
		queryCond += "AND prefect = $" + strconv.Itoa(pCount) + " "
		bind = append(bind, prefect)
	}

	// データの取得
	pCount++
	limitIdx := pCount
	pCount++
	offsetIdx := pCount
	queryData := "SELECT user_id, user_name, prefect, gender FROM users WHERE user_id != $1 " + queryCond + "LIMIT $" + strconv.Itoa(limitIdx) + " OFFSET $" + strconv.Itoa(offsetIdx)
	bindData := append(bind, limit, offset)

	users := []*model.UserWithoutPass{}
	err := ui.db.Select(&users, queryData, bindData...)
	if err != nil {
		return nil, err
	}

	// カウントの取得
	queryTotal := "SELECT COUNT(*) FROM users WHERE user_id != $1 " + queryCond
	var count int
	err = ui.db.Get(&count, queryTotal, bind...)
	if err != nil {
		return nil, err
	}

	return &model.UserList{
		HasNext: count > len(users)+offset,
		Users:   &users,
	}, nil
}

// get users by roomid
func (ui *userInfra) GetRoomUsers(roomId string, userId string) (*[]string, error) {
	users := []string{}
	err := ui.db.Select(
		&users,
		"SELECT users.user_id FROM users INNER JOIN room_datas ON room_datas.user_id = users.user_id WHERE room_datas.room_id = $1 AND users.user_id != $2",
		roomId,
		userId,
	)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

// get user not included in the room
func (ui *userInfra) GetUserListByUsername(limit int, offset int, name string, userId string, list []string) (*model.UserList, error) {
	query := "SELECT user_id, user_name, prefect, gender FROM users WHERE user_id != $1 AND user_name LIKE $2 LIMIT $3 OFFSET $4"
	users_p := []*model.UserWithoutPass{}
	err := ui.db.Select(
		&users_p,
		query,
		userId,
		"%"+name+"%",
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	users := []*model.UserWithoutPass{}
	for _, user_info := range users_p {
		isIncluded := false
		for _, userid := range list {
			if user_info.UserId == userid {
				isIncluded = true
				break
			}
		}
		if !isIncluded {
			users = append(users, user_info)
		}
	}

	queryCount := "SELECT COUNT(*) FROM users WHERE user_id != $1 AND user_name LIKE $2"
	var count int
	err = ui.db.Get(
		&count,
		queryCount,
		userId,
		"%"+name+"%",
	)
	if err != nil {
		return nil, err
	}

	return &model.UserList{
		HasNext: count-len(list) > len(users)+offset,
		Users:   &users,
	}, nil
}