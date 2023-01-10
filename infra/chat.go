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
		"SELECT `chat`.`post`, `chat`.`chat_id`, `chat`.`post_user_id`, `chat`.`created_at`, `chat`.`destination_user_id`, `room`.`room_id`, `room`.`not_read` FROM `room_data` as `room` INNER JOIN (SELECT * FROM `chats` GROUP BY `room_id` ORDER BY `created_at` DESC LIMIT 1) as `chat` ON `room`.`room_id` = `chat`.`room_id` AND `room`.`user_id` = ? ORDER BY `chat`.`created_at` DESC LIMIT ? OFFSET ?",
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
		if m.RoomId == "0" {
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
	query, args, err := sqlx.In("SELECT * FROM `users` WHERE `user_id` IN (?)", ids)
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
		"SELECT COUNT(*) FROM `room_data` WHERE `user_id` = ? GROUP BY `user_id`",
		userId,
	)

	res := &model.ChatList{
		HasNext: count > len(chatLst)+offset,
		Chats:   &chatLst,
	}

	return res, nil
}
