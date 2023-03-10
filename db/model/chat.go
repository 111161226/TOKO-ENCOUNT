package model

import "time"

//struct for receiving posted chat talk
type MessageSimple struct {
	Post string `json:"post" db:"post"`
}

//struct for message
type Message struct {
	Post      string    `json:"post" db:"post"`
	ChatId    string    `json:"chatId" db:"chat_id"`
	UserId    string    `json:"postUserId" db:"post_user_id"`
	UserName  string    `json:"userName" db:"user_name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

//struct for messages
type MessageList struct {
	HasNext  bool        `json:"hasNext"`
	Messages *[]*Message `json:"messages"`
}

//struct for return posted talk to chat
type ChatData struct {
	RoomId          string  `json:"roomId"`
	Name            string  `json:"name"`
	LatestMessage   Message `json:"latestMessage"`
	NewMessageCount int     `json:"newMessageCount"`
}

//struct for private chat talks
type ChatList struct {
	HasNext bool         `json:"hasNext"`
	Chats   *[]*ChatData `json:"chats"`
}

//struct for user's own room
type ChatUser struct {
	RoomId string `json:"roomId" db:"room_id"`
	UserId string `json:"userId" db:"user_id"`
}

//struct for user's all own room
type ChatUserList struct {
	ChatUsers *[]*ChatUser `json:"chatusers"`
}

//struct for room name
type RoomInfo struct {
	RoomId string `json:"roomId" db:"room_id"`
	RoomName   string `json:"roomName" db:"room_name"`
}
