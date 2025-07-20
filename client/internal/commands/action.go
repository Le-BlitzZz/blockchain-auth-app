package commands

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/Le-BlitzZz/blockchain-auth-app/client/internal/session"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

func triggerAction(action string) cli.ActionFunc {
	return func(c *cli.Context) error {
		sessionID, err := session.Start(action)
		if err != nil {
			return fmt.Errorf("could not start session: %w", err)
		}

		currentStatus := session.Started

		if err := browser.OpenURL("http://localhost:8080/" + sessionID); err != nil {
			return fmt.Errorf("could not open browser: %w", err)
		}

		ctx, cancel := signal.NotifyContext(c.Context, os.Interrupt)
		defer cancel()

		statusChan := session.Stream(ctx, sessionID)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case status, ok := <-statusChan:
					if !ok {
						return
					}

					switch status {
					case session.PendingWallet:
						if currentStatus == session.Started {
							fmt.Println("Please connect your wallet to continue.")
							currentStatus = session.PendingWallet
						}
					case session.WalletConnected:

					case session.Gone:
						if currentStatus == session.Started {
							fmt.Println("Please install MetaMask and retry the action.")
							return
						}
						fmt.Println("Gone.")
						return
					default:
						fmt.Printf("Unknown status received: %v\n", status)
					}
				case <-ctx.Done():
					return
				}
			}
		}()

		wg.Wait()

		return nil
	}
}
