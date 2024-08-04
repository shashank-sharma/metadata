package backend

import (
	"bytes"
	"encoding/json"
	"errors"
)

type LoginPayload struct {
	Identity string `json:"identity,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserRecord struct {
	Avatar          string `json:"avatar"`
	CollectionID    string `json:"collectionId"`
	CollectionName  string `json:"collectionName"`
	Created         string `json:"created"`
	Email           string `json:"email"`
	EmailVisibility bool   `json:"emailVisibility"`
	ID              string `json:"id"`
	Name            string `json:"name"`
	Updated         string `json:"updated"`
	Username        string `json:"username"`
	Verified        bool   `json:"verified"`
}

type LoginResponse struct {
	Record UserRecord `json:"record"`
	Token  string     `json:"token"`
}

func (bs *BackendService) Login(username, password string) (LoginResponse, error) {
	loginPayloadBuf := new(bytes.Buffer)
	json.NewEncoder(loginPayloadBuf).Encode(&LoginPayload{Identity: username, Password: password})
	req, _ := bs.Client.NewRequest("POST", "/api/collections/users/auth-with-password", loginPayloadBuf, map[string]string{})
	resp, err := bs.Client.Do(req)

	if err != nil {
		return LoginResponse{}, err
	}

	if resp.StatusCode != 200 {
		return LoginResponse{}, errors.New(resp.Status)
	}

	data := LoginResponse{}
	json.NewDecoder(resp.Body).Decode(&data)
	return data, nil
}
