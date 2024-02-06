package internal

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Chat struct {
	ID   string
	Name string
}

type ChatService struct {
	client *redis.Client
}

func NewChatService(client *redis.Client) *ChatService {
	return &ChatService{client}
}

func (s *ChatService) NewChat(ctx context.Context, name string) (*Chat, error) {
	id := uuid.NewString()
	key := fmt.Sprintf("chats:%s", id)
	err := s.client.Set(ctx, key, name, 0).Err()
	if err != nil {
		return nil, err
	}
	return &Chat{ID: id, Name: name}, err
}

func (s *ChatService) GetChat(ctx context.Context, id string) (*Chat, error) {
	key := fmt.Sprintf("chats:%s", id)
	name, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return &Chat{ID: id, Name: name}, nil
}
