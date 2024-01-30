package database

import (
	"github.com/google/uuid"
	"github.com/theghostmac/donorwallet/internal/models"
)

// CreateUser inserts a new user into the database.
func CreateUser(user *models.User) error {
	return DB.Create(user).Error
}

// GetUser retreives a user by their ID from the database.
func GetUser(userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := DB.Where("user_id =?", userID).First(&user).Error
	if err!= nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user in the database.
func UpdateUser(user *models.User) error {
	return DB.Save(user).Error
}

// DeleteUser removes a user from the database.
func DeleteUser(userID uuid.UUID) error {
	return DB.Where("user_id =?", userID).Delete(&models.User{}).Error
}