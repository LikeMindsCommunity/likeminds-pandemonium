package helpers

import (
	"fmt"
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/api/repository"
)

func GetSDKClientByApiKey(apiKey string) (*models.SDKClient, error) {
	sdkClient, err := repository.GetSDKClientByApiKey(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get sdk client, err=%s", err)
	}

	return sdkClient, nil
}
