package nixops

import (
	"backend/model"
	"fmt"

	"gorm.io/gorm"
)

func GetAllClients(db *gorm.DB) ([]model.Client, error) {
	var clients []model.Client
	if err := db.Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

func GetOneClient(db *gorm.DB, id uint) (model.Client, error) {
	var client model.Client
	if err := db.First(&client, id).Error; err != nil {
		return model.Client{}, fmt.Errorf("Client not found: %w", err)
	}
	return client, nil
}
