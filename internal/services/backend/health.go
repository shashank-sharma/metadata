package backend

import "errors"

func (bs *BackendService) HealthCheck() (bool, error) {
	req, _ := bs.Client.NewRequest("GET", "/api/health", nil, map[string]string{})
	resp, err := bs.Client.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode/100 != 2 {
		return false, errors.New(resp.Status)
	}

	return true, nil
}
