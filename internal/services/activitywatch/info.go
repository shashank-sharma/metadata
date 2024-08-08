package activitywatch

import (
	"encoding/json"
)

type AWInfo struct {
	Hostname string `json:"hostname"`
	Version  string `json:"version"`
	Testing  string `json:"testing"`
	DeviceId string `json:"device_id"`
}

func (as *AWService) FetchInfo() AWInfo {
	req, _ := as.Client.NewRequest("GET", "/api/0/info", nil, map[string]string{})
	resp, err := as.Client.Do(req)

	if err != nil {
		return AWInfo{}
	}

	if resp.StatusCode/100 != 2 {
		return AWInfo{}
	}

	data := AWInfo{}
	json.NewDecoder(resp.Body).Decode(&data)
	return data
}
