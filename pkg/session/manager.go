package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/theotruvelot/catchook/pkg/cache"
	"github.com/theotruvelot/catchook/storage/postgres/generated"
	"time"

	"github.com/redis/go-redis/v9"
)

type Manager interface {
	CreateSession(ctx context.Context, userID string, role string) (string, error)
	GetSession(ctx context.Context, sessionID string) (*Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
	RefreshSession(ctx context.Context, sessionID string) error
	ValidateSession(ctx context.Context, sessionID string) (*Session, error)
}

type Session struct {
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
	Role      generated.UserRole `json:"role"`
	CreatedAt time.Time          `json:"created_at"`
	ExpiresAt time.Time          `json:"expires_at"`
}

type sessionManager struct {
	redis    *redis.Client
	duration time.Duration
}

func NewManager(redis *redis.Client, duration time.Duration) Manager {
	return &sessionManager{
		redis:    redis,
		duration: duration,
	}
}

func (s *sessionManager) GetKey(sessionID string) string {
	return fmt.Sprintf(cache.KeyUserSession, sessionID)
}

func (s *sessionManager) CreateSession(ctx context.Context, userID string, role string) (string, error) {
	sessionID, err := s.generateSessionID()
	if err != nil {
		return "", fmt.Errorf("failed to generate session ID: %w", err)
	}

	now := time.Now()
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		Role:      generated.UserRole(role),
		CreatedAt: now,
		ExpiresAt: now.Add(s.duration),
	}

	key := s.GetKey(sessionID)
	data, err := json.Marshal(session)
	if err != nil {
		return "", fmt.Errorf("failed to marshal session: %w", err)
	}

	err = s.redis.Set(ctx, key, data, s.duration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store session: %w", err)
	}

	return sessionID, nil
}

func (s *sessionManager) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	key := s.GetKey(sessionID)
	data, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

func (s *sessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	key := s.GetKey(sessionID)
	err := s.redis.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func (s *sessionManager) RefreshSession(ctx context.Context, sessionID string) error {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now().Add(s.duration)

	key := s.GetKey(session.ID)
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	err = s.redis.Set(ctx, key, data, s.duration).Err()
	if err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	return nil
}

func (s *sessionManager) ValidateSession(ctx context.Context, sessionID string) (*Session, error) {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		err := s.DeleteSession(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("session expired")
	}

	return session, nil
}

func (s *sessionManager) generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
