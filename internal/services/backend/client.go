package backend

import (
	"time"

	"github.com/shashank-sharma/metadata/internal/services/apiclient"
)

type BackendService struct {
	Client *apiclient.APIClient
}

func NewBackendService(baseUrl string) *BackendService {
	apiClient := apiclient.New(baseUrl, 300*time.Second)
	return &BackendService{
		Client: apiClient,
	}
}
