package commands

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func pollStatus(session string) (string, error) {
	body, _ := json.Marshal(map[string]string{"session_id": session})
	resp, err := http.Post("http://localhost:8080/api/complete", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var res struct {
		Status string `json:"status"`
		Result string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.Status, nil
}

func pollStatusWithResult(session string) (status, result string, _ error) {
	body, _ := json.Marshal(map[string]string{"session_id": session})
	resp, err := http.Post("http://localhost:8080/api/complete", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	var res struct {
		Status string `json:"status"`
		Result string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", "", err
	}
	return res.Status, res.Result, nil
}
