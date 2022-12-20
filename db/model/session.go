package model

type Session struct {
	SessionId string `db:"session_id"`
	UserId    string `db:"user_id"`
}
