package infra

import (
	"github.com/cs-sysimpl/SakataKintoki/db/model"
	"github.com/cs-sysimpl/SakataKintoki/db/repository"
	"github.com/jmoiron/sqlx"
)

type chatInfra struct {
	db *sqlx.DB
}

func NewChatInfra(db *sqlx.DB) repository.ChatRepository {
	return &chatInfra{db: db}
}

func (ci *chatInfra) PostChat(chatId string, destinationId string, message *model.MessageSimple, post_user_id string, limit int, offset int) (*model.MessageList, error) {
	return nil, nil
}

func (ci *chatInfra) GetMessages(chatId string, limit int, offset int) (*model.MessageList, error) {
	return nil, nil
}

func (ci *chatInfra) CreateChat(destinationId string, post_user_id string) (*model.MessageList, error) {
	return nil, nil
}

func (ci *chatInfra) GetChatList(userId string, limit int, offset int) (*model.ChatList, error) {
	return nil, nil
}
