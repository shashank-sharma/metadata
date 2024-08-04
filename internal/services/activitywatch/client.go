package activitywatch

import (
	"time"

	"github.com/shashank-sharma/metadata/internal/services/apiclient"
)

type AWService struct {
	Client *apiclient.APIClient
}

func NewAWService(baseUrl string) *AWService {
	apiClient := apiclient.New(baseUrl, 60*time.Second)
	return &AWService{
		Client: apiClient,
	}
}
