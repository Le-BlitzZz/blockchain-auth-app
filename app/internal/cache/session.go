package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var Statuses = []string{
	SessionStatusStarted,
	SessionStatusPendingWallet,
	SessionStatusPendingSignature,
}

const (
	SessionStatusStarted          = "started"
	SessionStatusPendingWallet    = "pending_wallet"
	SessionStatusPendingSignature = "pending_signature"
)

const (
	WalletConnectExpiration = 15 * time.Second
	SignMessageExpiration   = 15 * time.Second
)

const sessionPrefix = "session:"

type Session struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Action string `json:"action,omitempty"`
}

func SessionKey(sessionId string) string {
	return fmt.Sprintf("%s%s", sessionPrefix, sessionId)
}

func NewSession(ctx context.Context, action string) (*Session, error) {
	sUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to create session ID: %w", err)
	}

	sId := sUUID.String()
	if sId == "" {
		return nil, fmt.Errorf("session ID is empty")
	}

	s := &Session{
		ID:     sId,
		Status: SessionStatusStarted,
		Action: action,
	}

	data, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session: %w", err)
	}

	return s, Redis().Set(ctx, SessionKey(sId), data, WalletConnectExpiration).Err()
}

func GetSession(ctx context.Context, sessionId string) (*Session, error) {
	data, err := Redis().Get(ctx, SessionKey(sessionId)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var s Session
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &s, nil
}

func UpdateSession(ctx context.Context, sessionId string, status string, timeout time.Duration) error {
	s, err := GetSession(ctx, sessionId)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	s.Status = status

	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	return Redis().Set(ctx, SessionKey(sessionId), data, timeout).Err()
}

func DeleteSession(ctx context.Context, sessionId string) error {
	return Redis().Del(ctx, SessionKey(sessionId)).Err()
}

func Sub(ctx context.Context) {
	Redis().PSubscribe(ctx, "").Channel(
		
	)
} 