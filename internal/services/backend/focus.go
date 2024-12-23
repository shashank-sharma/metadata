package backend

import (
	"bytes"
	"encoding/json"
	"errors"

	types "github.com/shashank-sharma/metadata/internal/types"
)

func (bs *BackendService) CreateTrackFocus(payload types.TrackFocusPayload) (bool, error) {
	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(&payload); err != nil {
		return false, err
	}

	req, err := bs.Client.NewRequest(
		"POST",
		"/api/focus/create",
		payloadBuf,
		map[string]string{},
	)
	if err != nil {
		return false, err
	}

	resp, err := bs.Client.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, errors.New(resp.Status)
	}

	return true, nil
}
