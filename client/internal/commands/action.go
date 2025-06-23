package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

const walletConnectTimeout = 20 * time.Second
const signMessageTimeout = 20 * time.Second

func triggerAction(action string) cli.ActionFunc {
	return func(c *cli.Context) error {
		startReq := map[string]string{"action": action}
		body, _ := json.Marshal(startReq)
		req, err := http.NewRequest("POST", "http://localhost:8080/api/start", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("failed to build start request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		// req.Header.Set("X-API-Key", config.APIKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to start session: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("start returned %d: %s", resp.StatusCode, string(bodyBytes))
		}

		var initResp struct {
			SessionID string `json:"session_id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&initResp); err != nil {
			return fmt.Errorf("bad start response: %w", err)
		}

		if err := browser.OpenURL("http://localhost:8080/" + initResp.SessionID); err != nil {
			return fmt.Errorf("could not open browser: %w", err)
		}

		var status string

		fmt.Println("→ Please connect your wallet in the browser (10s)")

		for range walletConnectTimeout {
			status, err := pollStatus(initResp.SessionID)
			if err != nil {
				return err
			}

			if status == "pending_signature" || status == "success" || status == "forbidden" {
				break
			}
			time.Sleep(1 * time.Second)
		}
		if status == "pending_wallet" {
			fmt.Println("Wallet connect timed out")
			return nil
		}

		fmt.Println("→ Please sign message in the browser (10s)")

		for range signMessageTimeout {
			status, result, err := pollStatusWithResult(initResp.SessionID)
			if err != nil {
				return err
			}

			switch status {
			case "pending_signature":
				time.Sleep(1 * time.Second) // wait before next poll
				continue
			case "success":
				fmt.Println(result) // should print "hello world" for "hello" action, "hello winner" for "vip", mint token for "mint"
				return nil
			case "expired":
				fmt.Println("session expired, please try again")
				return nil
			case "forbidden":
				fmt.Println("you do not hold a VIP pass, mint token")
				return nil
			default:
				println(err)
				return fmt.Errorf("unexpected status: %q", status)
			}
		}

		fmt.Println("timeout waiting for signature confirmation")
		return nil
	}
}
