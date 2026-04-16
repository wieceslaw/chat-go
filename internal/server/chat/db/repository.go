package db

import (
	"context"
	"database/sql"
	"log"
)

type UserRepository interface {
}

type ChatRepository interface {
	CreateChat(ctx context.Context)
	AddChatMember(ctx context.Context, chatId ChatId, userId UserId)
	RemoveChatMember(ctx context.Context, chatId ChatId, userId UserId)
}

type chatRepositoryImpl struct {
	db sql.DB
}

func NewChatRepository(ctx context.Context, connString string) (ChatRepository, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")

	// return &chatRepositoryImpl{
	// 	db: db,
	// }
	return nil, nil
}

func (cr *chatRepositoryImpl) CreateChat(ctx context.Context) {
	// cr.db.Exec()
}
