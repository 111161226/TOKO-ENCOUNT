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

//post message to chat
func (ci *chatInfra) PostChat(roomId string, message *model.MessageSimple, post_user_id string) (*model.Message, error) {
	ch, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	chatId := ch.String()
	//insert message into chats db
	_, err = ci.db.Exec(
		"INSERT INTO `chats` (`chat_id`, `room_id`, `post`, `post_user_id`) VALUES (?, ?, ?, ?)",
		chatId,
		roomId,
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
		"SELECT `post`, `chat_id`, `post_user_id`, `user_name`, `created_at` FROM `chats` INNER JOIN `users` ON `post_user_id` = `user_id` AND `chat_id` = ?",
		chatId,
	)

	//return posting message
	return &mess, nil
}

//get chat message by chat id
func (ci *chatInfra) GetMessages(roomId string, limit int, offset int) (*model.MessageList, error) {
	//get message info
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
	//get total of message 
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

func (ci *chatInfra) CreateChat(destinationId string, post_user_id string) (*model.ChatData, error) {
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
	m, err := ci.PostChat(roomId, &message, post_user_id)
	if err != nil {
		return nil, err
	}

	//return room info
	var name string
	err = ci.db.Get(
		&name,
		"SELECT `user_name` FROM `users` WHERE `user_id` = ?",
		destinationId,
	)
	if err != nil {
		return nil, err
	}

	//add room name
	var roomname = post_user_name + " " + name
	_, err = ci.db.Exec(
		"INSERT INTO `room_names` (`room_id`, `room_name`) VALUES (?, ?)",
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
		RoomId            string `db:"room_id"`
		NotRead           int    `db:"not_read"`
	}
	messages := []MessageDetail{}
	//get my joined rooms and those latest message
	err := ci.db.Select(
		&messages,
		"SELECT `chat`.`post`, `chat`.`chat_id`, `chat`.`post_user_id`, `chat`.`created_at`, `room`.`room_id`, `room`.`not_read` FROM (SELECT * FROM (SELECT `chat_id`, `room_id`, `post`, `post_user_id`, `created_at`, ROW_NUMBER() OVER(PARTITION BY `room_id` ORDER BY `created_at` DESC) AS `row_num` FROM `chats`) AS `c` WHERE `row_num` = 1) AS `chat` INNER JOIN (SELECT * FROM `room_datas` WHERE `user_id` = ?) AS `room` ON `room`.`room_id` = `chat`.`room_id` ORDER BY `chat`.`created_at` DESC LIMIT ? OFFSET ?",
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

	//select userId user name is needed
	idRoomMap := map[string][]string{}
	ids := []string{}
	for _, m := range messages {
		if m.RoomId == "0" { // case open chat
			if m.UserId != userId { //posting user is needed when the latest message is not sent by me
				ids = append(ids, m.UserId)
			}
			continue
		}

		if m.UserId == userId {
			//when the latest message is not sent by me
			var userids []string
			err = ci.db.Select(
				&userids,
				"SELECT `user_id` FROM `room_datas` WHERE `room_id` = ? AND `user_id` != ?",
				m.RoomId,
				m.UserId,
			)
			for _, id := range userids {
				idRoomMap[m.RoomId] = append(idRoomMap[m.RoomId], id)
				ids = append(ids, id)
			}
		} else {
			//when the latest message is sent by me
			ids = append(ids, m.UserId)
		}
	}
	ids = append(ids, userId) //want my username
	query, args, err := sqlx.In("SELECT `user_id`, `user_name`, `prefect`, `gender` FROM `users` WHERE `user_id` IN (?)", ids)
	if err != nil {
		return nil, err
	}
	users := []model.UserWithoutPass{}
	err = ci.db.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}

	//this dic has key in userId and value of username
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
			//when the latest message is sent by me
			tmp := ""
			for _, id := range idRoomMap[m.RoomId] {
				tmp += idNameMap[id] + ", "
			}
			name = tmp
		}
		c := &model.ChatData{
			RoomId:          m.RoomId,
			Name:            name,
			LatestMessage:   m.Message,
			NewMessageCount: m.NotRead,
		}
		chatLst = append(chatLst, c)
	}

	//get total of joined room number
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

//get users who join the designated chat ny roomId
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

//add new member to private chat
func (ci *chatInfra) AddPrivateChat(roomId string, did string) (*model.ChatData, error) {
	_, err := ci.db.Exec(
		"INSERT INTO `room_datas` (`room_id`, `user_id`) VALUES (?, ?)",
		roomId,
		did,
	)
	if err != nil {
		return nil, err
	}
	add_user_name := ""
	//create first message
	err = ci.db.Get(
		&add_user_name,
		"SELECT user_name FROM users WHERE user_id = ?",
		did,
	)
	if err != nil {
		return nil, err
	}
	mess := fmt.Sprintf("ユーザー%sが招待されました", add_user_name)
	message := model.MessageSimple{
		Post: mess,
	}
	//post first message
	m, err := ci.PostChat(roomId, &message, did)
	if err != nil {
		return nil, err
	}

	//return room info
	var name string
	err = ci.db.Get(
		&name,
		"SELECT `user_name` FROM `users` WHERE `user_id` = ?",
		did,
	)
	if err != nil {
		return nil, err
	}
	//get current username
	var curname string
	err = ci.db.Get(
		&curname,
		"SELECT `room_name` FROM `room_names` WHERE `room_id` = ?",
		roomId,
	)
	if err != nil {
		return nil, err
	}
	//update room name
	_, err = ci.db.Exec(
		"UPDATE `room_names` SET `room_name` = ? WHERE room_id = ?",
		curname+" " + name,
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

//get room name 
func (ci *chatInfra) GetRoomName(roomId string) (*model.RoomInfo, error) {
	//get room info
	room := model.RoomInfo{}
	err := ci.db.Get(
		&room,
		"SELECT `room_id`, `room_name` FROM `room_names` WHERE `room_id` = ? ",
		roomId,
	)
	return &room, err
}

//add new user to open chat 
func (ci *chatInfra) AddOpenChat(userId string) error {
	_, err := ci.db.Exec(
		"INSERT INTO `room_datas` (`room_id`, `user_id`) VALUES (0, ?)",
		userId,
	)
	return err
}

//not read message reset
func (ci *chatInfra) ResetNotRead(roomId string, userId string) error {
	_, err := ci.db.Exec(
		"UPDATE `room_datas` SET `not_read` = 0, `latest_access` = CURRENT_TIMESTAMP WHERE `room_id` = ? AND `user_id` = ?",
		roomId,
		userId,
	)
	return err
}

//not read message increment
func (ci *chatInfra) IncrementNotRead(roomId string, userId string) error {
	_, err := ci.db.Exec(
		"UPDATE `room_datas` SET `not_read` = `not_read` + 1 WHERE `room_id` = ? AND `user_id` != ?",
		roomId,
		userId,
	)
	return err
}
