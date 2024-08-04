package backend

import (
	"encoding/json"
	"errors"
)

type DevTokenResponse struct {
	Token string `json:"token"`
}

func (bs *BackendService) GetDevToken() (string, error) {
	req, _ := bs.Client.NewRequest("GET", "/api/token", nil, map[string]string{})
	resp, err := bs.Client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode/100 != 2 {
		return "", errors.New(resp.Status)
	}

	data := DevTokenResponse{}
	json.NewDecoder(resp.Body).Decode(&data)

	return data.Token, nil
}
