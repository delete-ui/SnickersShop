package handler

import (
	"SnickersShopPet1.0/internal/models"
	"SnickersShopPet1.0/internal/repository"
	"SnickersShopPet1.0/internal/validators"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

	if err := validators.AddSnickersValidate(&snickersInput); err != nil {
		m.repo.Zlog.Error("Error validating snickersInput", zap.Error(err), zap.String("location: ", op))
		http.Error(w, "Validation Error", http.StatusUnprocessableEntity)
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

func (m *SnickersHandler) AllSnickersGET(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handler.snickersHandler.AllSnickersGET"

	if r.Method != http.MethodGet {
		m.repo.Zlog.Error("Invalid request method", zap.String(" location: ", op))
		http.Error(w, "Invalid Request Method", http.StatusBadRequest)
		return
	}

	if r.Header.Get("content-type") != "application/json" {
		m.repo.Zlog.Error("Invalid content type", zap.String(" location: ", op))
		http.Error(w, "Invalid Content Type", http.StatusUnsupportedMediaType)
		return
	}

	var pagination models.Pagination
	if err := json.NewDecoder(r.Body).Decode(&pagination); err != nil {
		m.repo.Zlog.Error("Error decoding snickersInput", zap.Error(err), zap.String("location: ", op))
		http.Error(w, "Error Decoding Request", http.StatusInternalServerError)
		return
	}

	snickersSlice, err := m.repo.GetAll(&pagination)
	if err != nil {
		m.repo.Zlog.Error("Error getting snickers", zap.Error(err), zap.String("location: ", op))
		http.Error(w, "Error Getting Snickers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(snickersSlice); err != nil {
		m.repo.Zlog.Error("Error encoding snickers", zap.Error(err), zap.String("location: ", op))
		http.Error(w, "Error Encoding Response", http.StatusInternalServerError)
		return
	}

	m.repo.Zlog.Debug("SnickersHandler.AllSnickersGET successfully handled")

}

func (m *SnickersHandler) SnickersByIDGET(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handler.snickersHandler.SnickersByIDGET"

	if r.Method != http.MethodGet {
		m.repo.Zlog.Error("Invalid request method", zap.String(" location: ", op))
		http.Error(w, "Invalid Request Method", http.StatusBadRequest)
		return
	}

	idUUID := chi.URLParam(r, "id")

	snickers, err := m.repo.GetByID(idUUID)
	if err != nil {
		m.repo.Zlog.Error("Error getting snickers", zap.Error(err), zap.String("location: ", op))
		http.Error(w, "Error Getting Snickers", http.StatusInternalServerError)
		return
	}

	if snickers.ID == uuid.Nil {
		m.repo.Zlog.Debug("Shoes not found", zap.String(" location: ", op))
		http.Error(w, "Snickers Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(snickers); err != nil {
		m.repo.Zlog.Error("Error encoding snickers", zap.Error(err), zap.String("location: ", op))
		http.Error(w, "Error Encoding Response", http.StatusInternalServerError)
		return
	}

	m.repo.Zlog.Debug("SnickersHandler.SnickersByIDGET successfully handled")

}

func (m *SnickersHandler) SnickersByCostGET(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handler.snickersHandler.SnickersByCostGET"

	if r.Method != http.MethodGet {
		m.repo.Zlog.Error("Invalid request method", zap.String(" location: ", op))
		http.Error(w, "Invalid Request Method", http.StatusBadRequest)
		return
	}

	if r.Header.Get("content-type") != "application/json" {
		m.repo.Zlog.Error("Invalid content type", zap.String(" location: ", op))
		http.Error(w, "Invalid Content Type", http.StatusUnsupportedMediaType)
		return
	}

	var cost models.CostRange

	if err := json.NewDecoder(r.Body).Decode(&cost); err != nil {
		m.repo.Zlog.Error("Error decoding snickersInput", zap.Error(err), zap.String("location: ", op))
		http.Error(w, "Error Decoding Request Body", http.StatusInternalServerError)
		return
	}

	snickers, err := m.repo.GetByCost(&cost)
	if err != nil {
		m.repo.Zlog.Error("Error getting snickers", zap.Error(err), zap.String("location: ", op))
		http.Error(w, "Error Getting Snickers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(snickers); err != nil {
		m.repo.Zlog.Error("Error encoding snickers", zap.Error(err), zap.String("location: ", op))
		http.Error(w, "Error Encoding Response", http.StatusInternalServerError)
		return
	}

	m.repo.Zlog.Debug("SnickersHandler.SnickersByCostGET successfully handled")

}
