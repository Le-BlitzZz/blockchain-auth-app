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
							fmt.Println("Please connect your wallet to continue within 30 seconds.")
							currentStatus = session.PendingWallet
						}
					case session.PendingSignature:
						if currentStatus == session.Started || currentStatus == session.PendingWallet {
							fmt.Println("Please sign the message to continue within 30 seconds.")
							currentStatus = session.PendingSignature
						}
					case session.Verified:
						result, err := session.Finish(sessionID)
						if err != nil {
							fmt.Printf("Error finishing session: %v\n", err)
						} else {
							fmt.Println(result)
							return
						}
					case session.DeclinedSignature:
						currentStatus = session.DeclinedSignature
					case session.Gone:
						if currentStatus == session.Started {
							fmt.Println("Please install MetaMask and retry the action.")
							return
						}

						if currentStatus == session.PendingWallet {
							fmt.Println("Wallet connection timed out. Please retry the action.")
							return
						}

						if currentStatus == session.DeclinedSignature {
							fmt.Println("Signature request declined. Please retry the action.")
							return
						}

						if currentStatus == session.PendingSignature {
							fmt.Println("Signature request timed out. Please retry the action.")
							return
						}

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
