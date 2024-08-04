package backend

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Error struct {
	Error string `json:"error,omitempty"`
}

type TrackStatusPayload struct {
	Error
	UserId    string `json:"userid,omitempty"`
	Token     string `json:"token,omitempty"`
	ProductId string `json:"productid,omitempty"`
}

type TrackDevicePayload struct {
	Name     string `json:"name,omitempty"`
	HostName string `json:"hostname,omitempty"`
	Os       string `json:"os,omitempty"`
	Arch     string `json:"arch,omitempty"`
}

type TrackingDeviceResponse struct {
	ProductId string `json:"id"`
}

func (bs *BackendService) PingOnlineStatus(userId, token, productId string) (bool, error) {
	trackPayloadBuf := new(bytes.Buffer)
	json.NewEncoder(trackPayloadBuf).Encode(&TrackStatusPayload{UserId: userId, Token: token, ProductId: productId})
	req, _ := bs.Client.NewRequest("POST", "/api/track", trackPayloadBuf, map[string]string{})
	resp, err := bs.Client.Do(req)

	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, errors.New(resp.Status)
	}
	return true, nil
}

func (bs *BackendService) SetTrackingDevice(name, hostname, os, arch string) (string, error) {
	trackDevicePayloadBuf := new(bytes.Buffer)
	json.NewEncoder(trackDevicePayloadBuf).Encode(&TrackDevicePayload{Name: name, HostName: hostname, Os: os, Arch: arch})
	req, _ := bs.Client.NewRequest("POST", "/api/track/create", trackDevicePayloadBuf, map[string]string{})
	resp, err := bs.Client.Do(req)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	data := TrackingDeviceResponse{}
	json.NewDecoder(resp.Body).Decode(&data)
	return data.ProductId, nil
}
