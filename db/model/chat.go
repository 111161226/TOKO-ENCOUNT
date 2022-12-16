package model

import "time"

//struct for receiving posted chat talk
type ChatSimple struct {
	Post string `json:"post" db:"post"`
}

//struct for return posted talk to openchat
type Chat struct {
	Post              string    `json:"post" db:"post"`
	ChatId            string    `json:"chatId" db:"chat_id"`
	UserId            string    `json:"userId" db:"post_user_id"`
	DestinationUserId string    `json:"destinationUserId" db:"destination_user_id"`
	UserName          string    `json:"userName" db:"user_name"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
}

//struct for private chat talks
type ChatList struct {
	HasNext bool     `json:"hasNext"`
	Chats   *[]*Chat `json:"chats"`
}
