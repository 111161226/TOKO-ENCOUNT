package repository

import "github.com/cs-sysimpl/SakataKintoki/db/model"

type ChatRepository interface {
	PostChat(chatId string, destinationId string, message *model.MessageSimple, post_user_id string, limit int, offset int) (*model.MessageList, error)
	GetMessages(chatId string, limit int, offset int) (*model.MessageList, error)
	CreateChat(destinationId string, post_user_id string) (*model.MessageList, error)
}