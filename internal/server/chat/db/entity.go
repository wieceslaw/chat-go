package db

import (
	"time"
)

type (
	ChatId string
	UserId int64
)

type NewChatEntity struct {
	Id        ChatId
	Name      string
	OwnerId   UserId
	CreatedAt time.Time
}
