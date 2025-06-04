package handler

import (
	"SnickersShopPet1.0/internal/models"
	"SnickersShopPet1.0/internal/repository"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type SnickersHandler struct {
	repo *repository.SnickersRepository
}

func NewSnickersHandler(repo *repository.SnickersRepository) *SnickersHandler {
	return &SnickersHandler{repo: repo}
}

func (m *SnickersHandler) AddSnickersPOST(w http.ResponseWriter, r *http.Request) {

	const op = "internal.handler.snickersHandler.AddSnickersPOST"

	if r.Method != http.MethodPost {
		m.repo.Zlog.Error("Invalid request method", zap.String("location: ", op))
		http.Error(w, "Status Bad Request", http.StatusBadRequest)
		return
	}

	if r.Header.Get("content-type") != "application/json" {
		m.repo.Zlog.Error("Invalid content type", zap.String(" location: ", op))
		http.Error(w, "Invalid Content Type", http.StatusUnsupportedMediaType)
		return
	}

	var snickersInput models.SnickerInput

	if err := json.NewDecoder(r.Body).Decode(&snickersInput); err != nil {
		m.repo.Zlog.Error("Error decoding snickersInput", zap.Error(err), zap.String("location: ", op))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	snickers, err := m.repo.Add(&snickersInput)
	if err != nil {
		m.repo.Zlog.Error("Error adding snickers", zap.Error(err), zap.String("location: ", op))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(snickers); err != nil {
		m.repo.Zlog.Error("Error encoding snickers", zap.Error(err), zap.String("location: ", op))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	m.repo.Zlog.Debug("SnickersHandler.AddSnickersPOST successfully")

}
