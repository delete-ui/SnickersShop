package repository

import (
	"SnickersShopPet1.0/internal/models"
	"errors"
	"go.uber.org/zap"
)

type UserRepository struct {
	db   DBExecutor
	Zlog *zap.Logger
}

func NewUserRepository(db DBExecutor, logger *zap.Logger) *UserRepository {
	return &UserRepository{db: db, Zlog: logger}
}

func (h *UserRepository) AddUser(userInput *models.UserInput) (*models.User, error) {
	const op = "internal.userRepository.AddUser"

	err := h.db.QueryRow("SELECT * FROM products WHERE username = $1", userInput.Username)
	if err == nil {
		h.Zlog.Error("User already exists", zap.String(" location: ", op))
		return nil, errors.New("User already exists")
	}

	query := `INSERT INTO users (username, password)
	VALUES ($1, $2)
	RETURNING id, username, password;`

	row := h.db.QueryRow(query, userInput.Username, userInput.Password)

	var user models.User

	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		h.Zlog.Error("Error adding user", zap.Error(err), zap.String("location: ", op))
		return nil, err
	}

	return &user, nil

}

func (h *UserRepository) LoginUser(userInput *models.UserInput) (*models.UserLogInResponse, error) {
	const op = "internal.userRepository.LoginUser"

	query := `SELECT * FROM users WHERE username=$1;`

	row := h.db.QueryRow(query, userInput.Username)

	var user models.User

	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		h.Zlog.Error("Error getting user", zap.Error(err), zap.String("location: ", op))
		return nil, err
	}

	if user.Password != userInput.Password {
		return &models.UserLogInResponse{Username: userInput.Username, Status: "WRONG PASSWORD OR USERNAME"}, errors.New("wrong password")
	}

	return &models.UserLogInResponse{ID: user.ID, Username: user.Username, Status: "SUCCESSFULLY"}, nil

}
