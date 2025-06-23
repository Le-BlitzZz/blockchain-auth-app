package request

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var Statuses = []string{
	RequestStatusStarted,
	RequestStatusPendingWallet,
	RequestStatusPendingSignature,
}

const (
	RequestStatusStarted          = "started"
	RequestStatusPendingWallet    = "pending_wallet"
	RequestStatusPendingSignature = "pending_signature"
)

const (
	WalletConnectExpiration = 15 * time.Second
	SignMessageExpiration   = 15 * time.Second
)

const requestPrefix = "session:"

type Request struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Action string `json:"action,omitempty"`
}

func Key(requestId string) string {
	return fmt.Sprintf("%s%s", requestPrefix, requestId)
}

func NewRequest(ctx context.Context, action string) (*Request, error) {
	sUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to create request ID: %w", err)
	}

	sId := sUUID.String()
	if sId == "" {
		return nil, fmt.Errorf("request ID is empty")
	}

	s := &Request{
		ID:     sId,
		Status: RequestStatusStarted,
		Action: action,
	}

	data, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	return s, Redis().Set(ctx, Key(sId), data, WalletConnectExpiration).Err()
}

func GetRequest(ctx context.Context, requestId string) (*Request, error) {
	data, err := Redis().Get(ctx, Key(requestId)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get request: %w", err)
	}

	var s Request
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	return &s, nil
}

func UpdateRequest(ctx context.Context, requestId string, status string, timeout time.Duration) error {
	s, err := GetRequest(ctx, requestId)
	if err != nil {
		return fmt.Errorf("failed to get request: %w", err)
	}

	s.Status = status

	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	return Redis().Set(ctx, Key(requestId), data, timeout).Err()
}

func DeleteRequest(ctx context.Context, requestId string) error {
	return Redis().Del(ctx, Key(requestId)).Err()
}
