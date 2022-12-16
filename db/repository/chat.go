package repository

import "github.com/cs-sysimpl/SakataKintoki/db/model"

type ChatRepository interface {
	PostPchat(userId string, pchat *model.ChatSimple) (*model.ChatList, error)
}