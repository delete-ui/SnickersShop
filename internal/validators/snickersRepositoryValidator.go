package validators

import (
	"SnickersShopPet1.0/internal/models"
	"errors"
)

func AddSnickersValidate(snickers *models.SnickerInput) error {

	if snickers.Title == "" {
		return errors.New("Title can`t be empty")
	}
	if snickers.Cost <= 0 {
		return errors.New("Cost can`t be less than 0")
	}

	return nil

}
