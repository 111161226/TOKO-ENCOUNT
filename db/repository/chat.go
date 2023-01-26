package repository

import "github.com/111161226/TOKO-ENCOUNT/db/model"

type ChatRepository interface {
	PostChat(roomId string, destinationId string, message *model.MessageSimple, post_user_id string) (*model.Message, error)
	GetMessages(roomId string, limit int, offset int) (*model.MessageList, error)
	CreateChat(destinationId string, post_user_id string) (*model.ChatData, error)
	GetChatList(userId string, limit int, offset int) (*model.ChatList, error)
	GetChatByRoomId(roomId string) (*model.ChatUserList, error)
	AddOpenChat(userId string) error
	ResetNotRead(roomId string, userId string) error
	/** roomId に参加している userId 以外の全ユーザーの未読数をインクリメントする */
	IncrementNotRead(roomId string, userId string) error
}
