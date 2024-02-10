package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Chat struct {
	ID   string
	Name string
	TTL  time.Duration
}

func (c *Chat) ExpiresAt() string {
	expiration := time.Now().Add(c.TTL)
	return fmt.Sprintf("%s %d, %d:%d", expiration.Month().String(), expiration.Day(), expiration.Hour(), expiration.Minute())
}

type ChatService struct {
	client *redis.Client
}

const (
	ChatKey    = "chats:%s"
)

func NewChatService(client *redis.Client) *ChatService {
	return &ChatService{client}
}

func (s *ChatService) NewChat(ctx context.Context, name string) (*Chat, error) {
	id := uuid.NewString()
	key := fmt.Sprintf(ChatKey, id)
	ttl := 24 * time.Hour
	err := s.client.Set(ctx, key, name, ttl).Err()
	if err != nil {
		return nil, err
	}
	return &Chat{ID: id, Name: name, TTL: ttl}, err
}

func (s *ChatService) GetChat(ctx context.Context, id string) (*Chat, error) {
	key := fmt.Sprintf(ChatKey, id)
	name, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	ttl, err := s.client.TTL(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return &Chat{ID: id, Name: name, TTL: ttl}, nil
}

