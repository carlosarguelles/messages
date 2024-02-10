package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
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

type ChatMessage struct {
	Username string     `json:"username"`
	Content  string     `json:"content"`
	SentAt   *time.Time `json:"sentAt"`
}

type ChatService struct {
	client *redis.Client
}

const (
	ChatKey    = "chats:%s"
	MessageKey = ChatKey + ":messages:%s"
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

func (s *ChatService) NewMessage(ctx context.Context, chatID, username, content string) (*Message, error) {
	ttl, err := s.client.TTL(ctx, fmt.Sprintf(ChatKey, chatID)).Result()
	if err != nil {
		return nil, err
	}
	id := uuid.NewString()
	key := fmt.Sprintf(MessageKey, chatID, id)
	now := time.Now()
	value, err := json.Marshal(ChatMessage{username, content, &now})
	if err != nil {
		return nil, err
	}
	err = s.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return nil, err
	}
	return &Message{
		ChatID:  chatID,
		Content: MessageBubble(username, content, true),
	}, nil
}

func (s *ChatService) GetChatMessages(ctx context.Context, chatID string) ([]Message, error) {
	match := fmt.Sprintf(MessageKey, chatID, "*")
	keys := []string{}
	var cursor uint64
	for {
		k, c, err := s.client.Scan(ctx, cursor, match, 0).Result()
		if err != nil {
			continue
		}
		keys = append(keys, k...)
		cursor = c
		if cursor == 0 {
			break
		}
	}
	if len(keys) == 0 {
		return nil, nil
	}
	results, err := s.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	chatMessages := make([]ChatMessage, len(results))
	for i, result := range results {
		chatMessage := ChatMessage{}
		json.Unmarshal([]byte(result.(string)), &chatMessage)
		chatMessages[i] = chatMessage
	}
	slices.SortFunc[[]ChatMessage, ChatMessage](chatMessages, func(a, b ChatMessage) int {
		return a.SentAt.Compare(*b.SentAt)
	})
	messages := make([]Message, len(chatMessages))
	for i, cm := range chatMessages {
		messages[i] = Message{
			ChatID:  chatID,
			Content: MessageBubble(cm.Username, cm.Content, false),
		}
	}
	return messages, nil
}
