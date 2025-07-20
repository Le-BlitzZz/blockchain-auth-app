package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type StartRequestBody struct {
	Action string `json:"action"`
}

type StartResponseBody struct {
	SessionID string `json:"session_id"`
}

func Start(action string) (string, error) {
	reqBody, err := json.Marshal(StartRequestBody{Action: action})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(
		"http://localhost:8080/api/session",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to start session: %s", resp.Status)
	}

	var respBody StartResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", fmt.Errorf("bad start response: %w", err)
	}

	return respBody.SessionID, nil
}
