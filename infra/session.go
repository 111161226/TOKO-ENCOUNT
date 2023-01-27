package infra

import (
	"database/sql"
	"errors"

	"github.com/111161226/TOKO-ENCOUNT/db/model"
	"github.com/111161226/TOKO-ENCOUNT/db/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type sessionInfra struct {
	db *sqlx.DB
}

func NewSessionInfra(db *sqlx.DB) repository.SessionRepository {
	return &sessionInfra{db: db}
}

//create session
func (si *sessionInfra) CreateSession(userId string) (*model.Session, error) {
	uu, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	uuidStr := uu.String()

	_, err = si.db.Exec(
		"INSERT INTO `sessions` (`session_id`, `user_id`) VALUES (?, ?)",
		uuidStr,
		userId,
	)
	if err != nil {
		return nil, err
	}

	return &model.Session{
		SessionId: uuidStr,
		UserId:    userId,
	}, nil
}

//get seesion by session id
func (si *sessionInfra) GetSession(sessionId string) (*model.Session, error) {
	sess := model.Session{}
	err := si.db.Get(
		&sess,
		"SELECT * FROM `sessions` WHERE `session_id` = ?",
		sessionId,
	)
	if err != nil {
		return nil, err
	}

	return &sess, nil
}

//delete session by session id
func (si *sessionInfra) DeleteSessionBySessionId(sessionId string) error {
	_, err := si.db.Exec(
		"DELETE FROM `sessions` WHERE `session_id` = ?",
		sessionId,
	)
	return err
}

//delete session by user id
func (si *sessionInfra) DeleteSessionsByUserId(userId string) error {
	_, err := si.db.Exec(
		"DELETE FROM `sessions` WHERE `user_id` = ?",
		userId,
	)
	return err
}

//check session is present
func (si *sessionInfra) CheckSession(sessionId string) (*model.Session, error) {
	sess := model.Session{}
	err := si.db.Get(
		&sess,
		"SELECT * FROM `sessions` WHERE `session_id` = ?",
		sessionId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &sess, nil
}
