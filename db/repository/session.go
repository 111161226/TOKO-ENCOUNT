package repository

import "github.com/cs-sysimpl/SakataKintoki/db/model"

type SessionRepository interface {
	CreateSession(userId string) (*model.Session, error)
	GetSession(sessionId string) (*model.Session, error)
	DeleteSessionBySessionId(sessionId string) error
	DeleteSessionsByUserId(userId string) error
	CheckSession(sessionId string) (*model.Session, error)
}