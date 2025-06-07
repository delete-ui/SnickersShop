package handler

import (
	"SnickersShopPet1.0/internal/models"
	"SnickersShopPet1.0/internal/repository"
	"SnickersShopPet1.0/internal/validators"
	jwt2 "SnickersShopPet1.0/pkg/JWT"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) NewUserPOST(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handler.NewUserPOST"

	if r.Method != http.MethodPost {
		h.repo.Zlog.Error("Invalid request method", zap.String(" location:", op))
		http.Error(w, "Invalid Request Method", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		h.repo.Zlog.Error("Invalid request content-type", zap.String(" location:", op))
		http.Error(w, "Invalid Request Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	var userInput models.UserInput

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		h.repo.Zlog.Error("Error while decoding request body", zap.String(" location:", op), zap.Error(err))
		http.Error(w, "Error Decoding Request Body", http.StatusInternalServerError)
		return
	}

	if err := validators.ValidateAddUser(&userInput); err != nil {
		h.repo.Zlog.Error("Error while validating input", zap.String(" location:", op), zap.Error(err))
		http.Error(w, "Validation Error", http.StatusUnprocessableEntity)
		return
	}

	user, err := h.repo.AddUser(&userInput)
	if err != nil {
		h.repo.Zlog.Error("Error while adding user", zap.String(" location:", op), zap.Error(err))
		http.Error(w, "Error Creating User", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.repo.Zlog.Error("Error while encoding response", zap.String(" location:", op), zap.Error(err))
		http.Error(w, "Error Encoding Server Response", http.StatusInternalServerError)
		return
	}

	h.repo.Zlog.Debug("Method NewUserPOST successfully handled")

}

func (h *UserHandler) LogInPOST(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handler.LogInPOST"

	if r.Method != http.MethodPost {
		h.repo.Zlog.Error("Invalid request method", zap.String(" location:", op))
		http.Error(w, "Invalid Request Method", http.StatusBadRequest)
		return
	}

	if r.Header.Get("content-type") != "application/json" {
		h.repo.Zlog.Error("Invalid request content-type", zap.String(" location:", op))
		http.Error(w, "Invalid Request Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	var userInput models.UserInput

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		h.repo.Zlog.Error("Error while decoding request body", zap.String(" location:", op), zap.Error(err))
		http.Error(w, "Error Decoding Request Body", http.StatusInternalServerError)
		return
	}

	if err := validators.ValidateAddUser(&userInput); err != nil {
		h.repo.Zlog.Error("Error while validating input", zap.String(" location:", op), zap.Error(err))
		http.Error(w, "Validation Error", http.StatusUnprocessableEntity)
		return
	}

	userResponse, err := h.repo.LoginUser(&userInput)
	if err != nil {
		h.repo.Zlog.Error("Error while logging in", zap.String(" location:", op), zap.Error(err))
		http.Error(w, "Error Logging In", http.StatusInternalServerError)
		return
	}

	token, err := jwt2.GenerateToken(userResponse.ID)
	if err != nil {
		h.repo.Zlog.Error("Failed to generate token", zap.Error(err))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userResponse); err != nil {
		h.repo.Zlog.Error("Error while encoding response", zap.String(" location:", op), zap.Error(err))
		http.Error(w, "Error Encoding Server Response", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})

	h.repo.Zlog.Debug("Method LogInPOST successfully handled")

}
