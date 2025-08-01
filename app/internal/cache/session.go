package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	WalletConnectExpiration = 15 * time.Second
	SignMessageExpiration   = 15 * time.Second
)

const (
	SessionStatusStarted          = "started"
	SessionStatusPendingWallet    = "pending_wallet"
	SessionStatusPendingSignature = "pending_signature"
	SessionStatusVerified         = "verified"
)

type Session struct {
	ID        string  `json:"id" redis:"id"`
	Status    string  `json:"status" redis:"status"`
	Action    string  `json:"action" redis:"action"`
	Wallet    *string `json:"wallet" redis:"wallet"`
	Nonce     *string `json:"nonce,omitempty" redis:"nonce"`
	Message   *string `json:"message,omitempty" redis:"message"`
	Signature *string `json:"signature,omitempty" redis:"signature"`
	Result    *string `json:"result,omitempty" redis:"result"`
}

func NewSession(ctx context.Context, action string) *Session {
	s := &Session{
		Status: SessionStatusStarted,
		Action: action,
	}

	return s
}

func GetSession(ctx context.Context, sId string) (*Session, error) {
	exists, err := Redis().Exists(ctx, sId).Result()
	if err != nil {
		return nil, err
	}

	if exists == 0 {
		return nil, fmt.Errorf("session not found")
	}

	var s Session
	return &s, Redis().HGetAll(ctx, sId).Scan(&s)
}

func (s *Session) Create(ctx context.Context) error {
	sUUID, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to create session ID: %w", err)
	}

	sId := sUUID.String()
	s.ID = sId

	return Redis().HSet(ctx, s.ID, s).Err()
}

func (s *Session) Save(ctx context.Context, timeout time.Duration) error {
	err := Redis().HSet(ctx, s.ID, s).Err()
	if err != nil {
		return err
	}

	if timeout > 0 {
		return Redis().Expire(ctx, s.ID, timeout).Err()
	}

	return nil
}

func (s *Session) Delete(ctx context.Context) error {
	return Redis().Del(ctx, s.ID).Err()
}

func (s *Session) StreamSession(ctx context.Context) <-chan string {
	out := make(chan string)

	channel := "__keyspace@0__:" + s.ID
	pubsub := Redis().PSubscribe(ctx, channel)

	go func() {
		defer close(out)
		defer pubsub.Close()

		for {
			select {
			case msg, ok := <-pubsub.Channel():
				if !ok {
					return
				}
				switch msg.Payload {
				case "hset":
					updatedSess, err := GetSession(ctx, s.ID)
					if err != nil {
						log.Error("Failed to get session:", err)
						return
					}
					if updatedSess.Status != s.Status {
						out <- updatedSess.Status
					}
					*s = *updatedSess
				case "del", "expired":
					out <- "gone"
					return
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
