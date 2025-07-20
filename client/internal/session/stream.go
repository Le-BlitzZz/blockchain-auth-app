package session

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/r3labs/sse/v2"
)

func Stream(ctx context.Context, sessionID string) <-chan Status {
	client := sse.NewClient("http://localhost:8080/api/session/" + sessionID + "/stream")

	statusChan := make(chan Status, 1)

	go func() {
		defer close(statusChan)

		err := client.SubscribeWithContext(ctx, "", func(msg *sse.Event) {
			var status Status

			data := string(msg.Data)
			switch data {
			case "pending_wallet":
				status = PendingWallet
			case "wallet_connected":
				status = WalletConnected
			case "pending_signature":
				status = PendingSignature
			case "gone":
				status = Gone
			default:
				fmt.Printf("Unknown status: %s\n", data)
				return
			}

			select {
			case statusChan <- status:
			case <-ctx.Done():
				return
			}

		})

		if err != nil {
			log.Info("Failed to subscribe to session stream:", err)
			return
		}
	}()

	return statusChan
}
