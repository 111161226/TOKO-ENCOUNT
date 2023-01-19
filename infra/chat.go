package infra

import (
	"fmt"

	"github.com/cs-sysimpl/SakataKintoki/db/model"
	"github.com/cs-sysimpl/SakataKintoki/db/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
		"SELECT post, chat_id, post_user_id, user_name, created_at FROM chats INNER JOIN users ON post_user_id = user_id WHERE room_id = ? AND post_user_id = ? ORDER BY `created_at` DESC",
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
		"SELECT post, chat_id, post_user_id, user_name, created_at FROM chats INNER JOIN users ON post_user_id = user_id WHERE room_id = ? ORDER BY `created_at` DESC LIMIT ? OFFSET ?",
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
		HasNext:  count > len(mess),
	}, nil
}

func (ci *chatInfra) CreateChat(destinationId string, post_user_id string) (*model.Message, error) {
	post_user_name := ""
	//create first message
	err := ci.db.Get(
		&post_user_name,
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
	return ci.PostChat(roomId, destinationId, &message, post_user_id)
}

func (ci *chatInfra) GetChatList(userId string, limit int, offset int) (*model.ChatList, error) {
	type MessageDetail struct {
		model.Message
		RoomId            string `db:"room_id"`
		DestinationUserId string `db:"destination_user_id"`
		NotRead           int    `db:"not_read"`
	}
	messages := []MessageDetail{}
	// 自分が参加しているルームとその最新メッセージを取得
	err := ci.db.Select(
		&messages,
		"SELECT `chat`.`post`, `chat`.`chat_id`, `chat`.`post_user_id`, `chat`.`created_at`, `chat`.`destination_user_id`, `room`.`room_id`, `room`.`not_read` FROM `room_datas` as `room` INNER JOIN (SELECT * FROM `chats` GROUP BY `room_id` ORDER BY `created_at` DESC LIMIT 1) as `chat` ON `room`.`room_id` = `chat`.`room_id` AND `room`.`user_id` = ? ORDER BY `chat`.`created_at` DESC LIMIT ? OFFSET ?",
		userId,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		res := &model.ChatList{
			HasNext: false,
			Chats:   &[]*model.ChatData{},
		}
		return res, nil
	}

	// ユーザー名情報が必要なユーザーIDを列挙
	ids := []string{}
	for _, m := range messages {
		if m.RoomId == "0" { // 全体チャットの場合
			if m.UserId != userId { // 自分が送っていないメッセージの場合、送信者の名前が必要
				ids = append(ids, m.UserId)
			}
			continue
		}

		if m.UserId == userId {
			// 最新メッセージが自分が送ったものの場合
			ids = append(ids, m.DestinationUserId)
		} else {
			// 最新メッセージが相手が送ったものの場合
			ids = append(ids, m.UserId)
		}
	}
	ids = append(ids, userId) // 自分のユーザー名も欲しい
	query, args, err := sqlx.In("SELECT `user_id`, `user_name`, `prefect`, `gender` FROM `users` WHERE `user_id` IN (?)", ids)
	if err != nil {
		return nil, err
	}
	users := []model.UserWithoutPass{}
	err = ci.db.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}

	// ユーザーIDをキー、ユーザー名を値に持つ
	idNameMap := map[string]string{}
	for _, u := range users {
		idNameMap[u.UserId] = u.UserName
	}

	chatLst := []*model.ChatData{}
	for _, m := range messages {
		m.UserName = idNameMap[m.UserId]
		name := idNameMap[m.UserId]
		if m.RoomId == "0" {
			name = "全体チャット"
		} else if m.UserId == userId {
			// 最新メッセージが自分が送ったものの場合
			name = idNameMap[m.DestinationUserId]

		}
		c := &model.ChatData{
			RoomId:          m.RoomId,
			Name:            name,
			LatestMessage:   m.Message,
			NewMessageCount: m.NotRead,
		}
		chatLst = append(chatLst, c)
	}

	// 自分が参加しているルームの全数を取得
	count := 0
	err = ci.db.Get(
		&count,
		"SELECT COUNT(*) FROM `room_datas` WHERE `user_id` = ? GROUP BY `user_id`",
		userId,
	)

	res := &model.ChatList{
		HasNext: count > len(chatLst)+offset,
		Chats:   &chatLst,
	}

	return res, nil
}

func (ci *chatInfra) GetChatByRoomId(roomId string) (*model.ChatUserList, error) {
	users := []*model.ChatUser{}
	err := ci.db.Select(
		&users,
		"SELECT room_id, user_id FROM room_datas WHERE room_id = ?",
		roomId,
	)
	if err != nil {
		return nil, err
	}
	return &model.ChatUserList{
		ChatUsers: &users,
	}, nil
}
