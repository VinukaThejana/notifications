package lib

import "github.com/VinukaThejana/notifications/internal/models"

// FindUserByID finds the user from the givern user array from the user id
func FindUserByID(id int, users []models.User) *models.User {
	for _, user := range users {
		if user.ID == id {
			return &user
		}
	}

	return nil
}
