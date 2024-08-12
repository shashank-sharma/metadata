package backend

import (
	"bytes"
	"encoding/json"
	"errors"

	types "github.com/shashank-sharma/metadata/internal/types"
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

type EventListPayload struct {
	DeviceId string          `json:"device_id"`
	TaskName string          `json:"task_name"`
	Events   []types.AWEvent `json:"events"`
}

type TrackingDeviceResponse struct {
	ProductId string `json:"id"`
}

type EventSyncResponse struct {
	CreateCount int64 `json:"create_count"`
	FailedCount int64 `json:"failed_count"`
	SkipCount   int64 `json:"skip_count"`
	ForceCheck  bool  `json:"force_check"`
}

func (bs *BackendService) SyncEventData(deviceId string, taskName string, events []types.AWEvent) (EventSyncResponse, error) {
	eventPayloadBuf := new(bytes.Buffer)
	json.NewEncoder(eventPayloadBuf).Encode(&EventListPayload{DeviceId: deviceId, TaskName: taskName, Events: events})
	req, _ := bs.Client.NewRequest("POST", "/api/sync/create", eventPayloadBuf, map[string]string{})
	resp, err := bs.Client.Do(req)

	if err != nil {
		return EventSyncResponse{}, err
	}

	if resp.StatusCode != 200 {
		return EventSyncResponse{}, errors.New(resp.Status)
	}

	data := EventSyncResponse{}
	json.NewDecoder(resp.Body).Decode(&data)
	return data, nil
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
