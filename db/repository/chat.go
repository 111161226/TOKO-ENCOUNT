package repository

import "github.com/cs-sysimpl/SakataKintoki/db/model"

type ChatRepository interface {
	PostChat(roomId string, destinationId string, message *model.MessageSimple, post_user_id string) (*model.Message, error)
	GetMessages(roomId string, limit int, offset int) (*model.MessageList, error)
	CreateChat(destinationId string, post_user_id string) (*model.Message, error)
	GetChatList(userId string, limit int, offset int) (*model.ChatList, error)
}
