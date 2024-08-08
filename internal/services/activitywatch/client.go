package activitywatch

import (
	"time"

	"github.com/shashank-sharma/metadata/internal/services/apiclient"
)

type AWService struct {
	Client *apiclient.APIClient
	AWInfo AWInfo
}

func NewAWService(baseUrl string) *AWService {
	apiClient := apiclient.New(baseUrl, 60*time.Second)
	awService := &AWService{
		Client: apiClient,
	}
	awService.AWInfo = awService.FetchInfo()
	return awService
}
