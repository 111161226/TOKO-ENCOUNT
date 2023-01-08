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

func (ci *chatInfra) PostChat(roomId string, destinationId string, message *model.MessageSimple, post_user_id string) (*model.Message, error) {
	ch, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	chatId := ch.String()
	//insert message into chats db
	_, err = ci.db.Exec(
		"INSERT INTO `chats` (`chat_id`, `room_id`, `destination_user_id`, `post`, `post_user_id`) VALUES (?, ?, ?, ?, ?)",
		chatId,
		roomId,
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
		"SELECT post, room_id, post_user_id, user_name, created_at FROM chats INNER JOIN users ON post_user_id = user_id WHERE room_id = ? AND post_user_id = ? ORDER BY `created_at` DESC",
		roomId,
		post_user_id,
	)

	//return posting message 
	return &mess, nil
}

func (ci *chatInfra) GetMessages(roomId string, limit int, offset int) (*model.MessageList, error) {
	mess := []*model.Message{}
	err := ci.db.Select(
		&mess,
		"SELECT post, room_id, post_user_id, user_name, created_at FROM chats INNER JOIN users ON post_user_id = user_id WHERE room_id = ? ORDER BY `created_at` DESC LIMIT ? OFFSET ?",
		roomId,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	count := 0
	err = ci.db.Get(
		&count,
		"SELECT COUNT(*) FROM `chats` WHERE `room_id` = ?",
		roomId,
	)
	if err != nil {
		return nil, err
	}
	return &model.MessageList{
		Messages: &mess,
		HasNext: count > len(mess),
	}, nil
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
	roomId := ch.String()
	//add room_data to db
	_, err = ci.db.Exec(
		"INSERT INTO `room_datas` (`room_id`, `user_id`) VALUES (?, ?), (?, ?)",
		roomId,
		post_user_id,
		roomId,
		destinationId,
	)
	if err != nil {
		return nil, err
	}
	//post first message
	return ci.PostChat(roomId, destinationId , &message, post_user_id) 
}

func (ci *chatInfra) GetChatList(userId string, limit int, offset int) (*model.ChatList, error) {
	return nil, nil
}

func (ci *chatInfra) GetChatByRoomId(roomId string) (*model.ChatUserList, error) {
	users := []*model.ChatUser{}
	err := ci.db.Select(
		&users,
		"SELECT `room_id`, `user_id` FROM room_datas WHRER room_id = ?",
		roomId,
	)
	if err != nil {
		return nil, err
	}
	return &model.ChatUserList{
		ChatUsers: &users,
	}, nil
}