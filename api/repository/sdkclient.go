package repository

import (
	"fmt"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/models"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/database"

	"gorm.io/gorm"
)

type SDKClientRepository struct {
	sdkClientDatabase *gorm.DB
}

func NewSDKClientRepositoryRepository() *SDKClientRepository {
	return &SDKClientRepository{
		sdkClientDatabase: database.Postgres,
	}
}

func GetSDKClientByApiKey(apiKey string) (*models.SDKClient, error) {
	sdkClient := &models.SDKClient{}

	filter := map[string]interface{}{
		"api_key": apiKey,
	}

	err := NewSDKClientRepositoryRepository().sdkClientDatabase.Where(filter).First(&sdkClient).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get sdk client, err=%s", err)
	}

	return sdkClient, nil
}
