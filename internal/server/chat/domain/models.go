package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	ChatEventType int
	UserId        int64
	ChatId        uuid.UUID
	MessageId     uuid.UUID
)

type ChatMessage struct {
	Id        MessageId
	Text      string
	Timestamp time.Time
	AuthorId  UserId
	ChatId    ChatId
}

type NewChat struct {
	Id        ChatId
	Name      string
	OwnerId   UserId
	CreatedAt time.Time
}

const (
	MessageEvent ChatEventType = iota
	ChatJoinEvent
	ChatLeaveEvent
)

type ChatEvent struct {
	Type      ChatEventType
	Message   string
	Timestamp time.Time
}
