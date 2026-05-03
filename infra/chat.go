package infra

import (
	"fmt"

	"github.com/111161226/TOKO-ENCOUNT/db/model"
	"github.com/111161226/TOKO-ENCOUNT/db/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type chatInfra struct {
	db *sqlx.DB
}

func NewChatInfra(db *sqlx.DB) repository.ChatRepository {
	return &chatInfra{db: db}
}

// post message to chat
func (ci *chatInfra) PostChat(roomId string, message *model.MessageSimple, post_user_id string) (*model.Message, error) {
	ch, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	chatId := ch.String()
	// insert message into chats db
	_, err = ci.db.Exec(
		"INSERT INTO chats (chat_id, room_id, post, post_user_id) VALUES ($1, $2, $3, $4)",
		chatId,
		roomId,
		message.Post,
		post_user_id,
	)
	if err != nil {
		return nil, err
	}
	// get posting message
	mess := model.Message{}
	err = ci.db.Get(
		&mess,
		"SELECT post, chat_id, post_user_id, user_name, created_at FROM chats INNER JOIN users ON post_user_id = user_id AND chat_id = $1",
		chatId,
	)

	return &mess, nil
}

// get chat message by chat id
func (ci *chatInfra) GetMessages(roomId string, limit int, offset int) (*model.MessageList, error) {
	mess := []*model.Message{}
	err := ci.db.Select(
		&mess,
		"SELECT post, chat_id, post_user_id, user_name, created_at FROM chats INNER JOIN users ON post_user_id = user_id WHERE room_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
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
		"SELECT COUNT(*) FROM chats WHERE room_id = $1",
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

func (ci *chatInfra) CreateChat(destinationId string, post_user_id string) (*model.ChatData, error) {
	post_user_name := ""
	err := ci.db.Get(
		&post_user_name,
		"SELECT user_name FROM users WHERE user_id = $1",
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
	// add room_data to db
	_, err = ci.db.Exec(
		"INSERT INTO room_datas (room_id, user_id) VALUES ($1, $2), ($3, $4)",
		roomId,
		post_user_id,
		roomId,
		destinationId,
	)
	if err != nil {
		return nil, err
	}

	m, err := ci.PostChat(roomId, &message, post_user_id)
	if err != nil {
		return nil, err
	}

	var name string
	err = ci.db.Get(
		&name,
		"SELECT user_name FROM users WHERE user_id = $1",
		destinationId,
	)
	if err != nil {
		return nil, err
	}

	var roomname = post_user_name + ", " + name
	_, err = ci.db.Exec(
		"INSERT INTO room_names (room_id, room_name) VALUES ($1, $2)",
		roomId,
		roomname,
	)
	if err != nil {
		return nil, err
	}

	return &model.ChatData{
		RoomId:          roomId,
		Name:            roomname,
		LatestMessage:   *m,
		NewMessageCount: 0,
	}, nil
}

func (ci *chatInfra) GetChatList(userId string, limit int, offset int) (*model.ChatList, error) {
	type MessageDetail struct {
		model.Message
		RoomId  string `db:"room_id"`
		NotRead int    `db:"not_read"`
	}
	messages := []MessageDetail{}
	//get my joined rooms and those latest message
	err := ci.db.Select(
		&messages,
		"SELECT chat.post, chat.chat_id, chat.post_user_id, chat.created_at, room.room_id, room.not_read FROM (SELECT * FROM (SELECT chat_id, room_id, post, post_user_id, created_at, ROW_NUMBER() OVER(PARTITION BY room_id ORDER BY created_at DESC) AS row_num FROM chats) AS c WHERE row_num = 1) AS chat INNER JOIN (SELECT * FROM room_datas WHERE user_id = $1) AS room ON room.room_id = chat.room_id ORDER BY chat.created_at DESC LIMIT $2 OFFSET $3",
		userId,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return &model.ChatList{
			HasNext: false,
			Chats:   &[]*model.ChatData{},
		}, nil
	}

	//get room names by roomId
	RoomidNameMap := map[string]string{}
	for _, m := range messages {
		if m.RoomId == "0" {
			RoomidNameMap[m.RoomId] = "全体チャット"
			continue
		}
		var roomname string
		err = ci.db.Get(
			&roomname,
			"SELECT room_name FROM room_names WHERE room_id = $1",
			m.RoomId,
		)
		if err != nil {
			return nil, err
		}
		RoomidNameMap[m.RoomId] = roomname
	}

	chatLst := []*model.ChatData{}
	for _, m := range messages {
		c := &model.ChatData{
			RoomId:          m.RoomId,
			Name:            RoomidNameMap[m.RoomId],
			LatestMessage:   m.Message,
			NewMessageCount: m.NotRead,
		}
		chatLst = append(chatLst, c)
	}

	count := 0
	// PostgreSQL は GROUP BY で指定したカラム以外を SELECT する際に制約があるが、COUNT(*) のみなら問題なし
	err = ci.db.Get(
		&count,
		"SELECT COUNT(*) FROM room_datas WHERE user_id = $1",
		userId,
	)

	return &model.ChatList{
		HasNext: count > len(chatLst)+offset,
		Chats:   &chatLst,
	}, nil
}

func (ci *chatInfra) UpdateRoomName(roomId string, name string, userId string) (*model.ChatData, error) {
	_, err := ci.db.Exec(
		"UPDATE room_names SET room_name = $1 WHERE room_id = $2 ",
		name,
		roomId,
	)
	if err != nil {
		return nil, err
	}
	post_user_name := ""
	err = ci.db.Get(
		&post_user_name,
		"SELECT user_name FROM users WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return nil, err
	}
	mess := fmt.Sprintf("ユーザー%sが名前を%sに変更しました", post_user_name, name)
	message := model.MessageSimple{
		Post: mess,
	}
	postedMessage, err := ci.PostChat(roomId, &message, userId)
	if err != nil {
		return nil, err
	}
	return &model.ChatData{
		RoomId:          roomId,
		Name:            name,
		LatestMessage:   *postedMessage,
		NewMessageCount: 0,
	}, nil
}

//get users who join the designated chat ny roomId
func (ci *chatInfra) GetChatByRoomId(roomId string) (*model.ChatUserList, error) {
	users := []*model.ChatUser{}
	err := ci.db.Select(
		&users,
		"SELECT room_id, user_id FROM room_datas WHERE room_id = $1",
		roomId,
	)
	if err != nil {
		return nil, err
	}
	return &model.ChatUserList{
		ChatUsers: &users,
	}, nil
}

//add new member to private chat
func (ci *chatInfra) AddPrivateChat(roomId string, did string, post_user_id string) (*model.ChatData, error) {
	_, err := ci.db.Exec(
		"INSERT INTO room_datas (room_id, user_id) VALUES ($1, $2)",
		roomId,
		did,
	)
	if err != nil {
		return nil, err
	}
	add_user_name := ""
	err = ci.db.Get(
		&add_user_name,
		"SELECT user_name FROM users WHERE user_id = $1",
		did,
	)
	if err != nil {
		return nil, err
	}
	mess := fmt.Sprintf("ユーザー%sが招待されました", add_user_name)
	message := model.MessageSimple{
		Post: mess,
	}
	m, err := ci.PostChat(roomId, &message, post_user_id)
	if err != nil {
		return nil, err
	}

	var name string
	err = ci.db.Get(
		&name,
		"SELECT user_name FROM users WHERE user_id = $1",
		did,
	)
	if err != nil {
		return nil, err
	}
	var curname string
	err = ci.db.Get(
		&curname,
		"SELECT room_name FROM room_names WHERE room_id = $1",
		roomId,
	)
	if err != nil {
		return nil, err
	}
	_, err = ci.db.Exec(
		"UPDATE room_names SET room_name = $1 WHERE room_id = $2",
		curname+", "+name,
		roomId,
	)
	if err != nil {
		return nil, err
	}

	return &model.ChatData{
		RoomId:          roomId,
		Name:            name,
		LatestMessage:   *m,
		NewMessageCount: 0,
	}, nil
}

func (ci *chatInfra) GetRoomName(roomId string) (*model.RoomInfo, error) {
	room := model.RoomInfo{}
	err := ci.db.Get(
		&room,
		"SELECT room_id, room_name FROM room_names WHERE room_id = $1 ",
		roomId,
	)
	return &room, err
}

func (ci *chatInfra) AddOpenChat(userId string) error {
	_, err := ci.db.Exec(
		"INSERT INTO room_datas (room_id, user_id) VALUES ('0', $1)",
		userId,
	)
	return err
}

func (ci *chatInfra) ResetNotRead(roomId string, userId string) error {
	_, err := ci.db.Exec(
		"UPDATE room_datas SET not_read = 0, latest_access = CURRENT_TIMESTAMP WHERE room_id = $1 AND user_id = $2",
		roomId,
		userId,
	)
	return err
}

func (ci *chatInfra) IncrementNotRead(roomId string, userId string) error {
	_, err := ci.db.Exec(
		"UPDATE room_datas SET not_read = not_read + 1 WHERE room_id = $1 AND user_id != $2",
		roomId,
		userId,
	)
	return err
}
