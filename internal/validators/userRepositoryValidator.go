package validators

import (
	"SnickersShopPet1.0/internal/models"
	"errors"
)

func ValidateAddUser(user *models.UserInput) error {

	if len(user.Username) < 3 {
		return errors.New("Username must be at least 3 characters long")
	}
	if len(user.Password) < 6 {
		return errors.New("Password must be at least 6 characters long")
	}

	return nil

}
