package domain

import (
	"context"

	"github.com/wieceslaw/chat-go/internal/server/chat/db"
)

type ChatService interface {
	CreateChat(ctx context.Context, chat NewChat)
	RemoveChat(ctx context.Context, chatId ChatId)

	JoinChat(ctx context.Context, chatId ChatId, userId UserId)
	LeaveChat(ctx context.Context, chatId ChatId, userId UserId)

	SendMessage(ctx context.Context, message ChatMessage) MessageId
	GetMessagesPage(ctx context.Context, offset int, count int) []ChatMessage
}

func NewChatService(repository db.ChatRepository) ChatService {
	// return &chatServiceImpl{
	// 	repository: repository,
	// }
	return nil
}

type chatServiceImpl struct {
	repository db.ChatRepository
}

func (cr *chatServiceImpl) CreateChat(ctx context.Context, chat NewChat) {
	// cr.repository.CreateChat()
}
