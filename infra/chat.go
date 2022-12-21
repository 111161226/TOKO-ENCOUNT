package infra

import (
	"fmt"

	"github.com/cs-sysimpl/SakataKintoki/db/model"
	"github.com/cs-sysimpl/SakataKintoki/db/repository"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
)

type chatInfra struct {
	db *sqlx.DB
}

func NewChatInfra(db *sqlx.DB) repository.ChatRepository {
	return &chatInfra{db: db}
}

func (ci *chatInfra) PostChat(chatId string, destinationId string, message *model.MessageSimple, post_user_id string) (*model.Message, error) {
	//insert message into chats db
	_, err := ci.db.Exec(
		"INSERT INTO `chats` (`chat_id`, `destination_user_id`, `post`, `post_user_id`) VALUES (?, ?, ?, ?)",
		chatId,
		destinationId,
		message.Post,
		post_user_id,
	)
	if err != nil {
		return nil, err
	}
	//get posting message
	mess := model.Message{}
	err = ci.db.Get(
		&mess,
		"SELECT post, chat_id, post_user_id, user_name, created_at FROM chats INNER JOIN users ON post_user_id = user_id WHERE chat_id = ? AND post_user_id = ? ORDER BY `created_at` DESC",
		chatId,
		post_user_id,
	)

	//return posting message 
	return &mess, nil
}

func (ci *chatInfra) GetMessages(chatId string, limit int, offset int) (*model.MessageList, error) {
	return nil, nil
}

func (ci *chatInfra) CreateChat(destinationId string, post_user_id string) (*model.Message, error) {
	post_user_name := ""
	//create first message
	err := ci.db.Get(
		&post_user_id,
		"SELECT user_name FROM users WHERE user_id = ?",
		post_user_id,
	)
	if err != nil {
		return nil, err
	}
	mess := fmt.Sprintf("ユーザー%sがチャットを始めました", post_user_name)
	message := model.MessageSimple{
		Post: mess,
	}
	ch, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	chatId := ch.String()
	//add room_data to db
	_, err = ci.db.Exec(
		"INSERT INTO `room_datas` (`room_id`, `user_id`) VALUES (?, ?)",
		chatId,
		post_user_id,
	)
	if err != nil {
		return nil, err
	}
	_, err = ci.db.Exec(
		"INSERT INTO `room_datas` (`room_id`, `user_id`) VALUES (?, ?)",
		chatId,
		destinationId,
	)
	if err != nil {
		return nil, err
	}
	//post first message
	return ci.PostChat(chatId, destinationId , &message, post_user_id) 
}


func (ci *chatInfra) GetChatList(userId string, limit int, offset int) (*model.ChatList, error) {
	return nil, nil
}
