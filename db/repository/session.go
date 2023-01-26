package repository

import "github.com/111161226/TOKO-ENCOUNT/db/model"

type SessionRepository interface {
	CreateSession(userId string) (*model.Session, error)
	GetSession(sessionId string) (*model.Session, error)
	DeleteSessionBySessionId(sessionId string) error
	DeleteSessionsByUserId(userId string) error
	CheckSession(sessionId string) (*model.Session, error)
}
