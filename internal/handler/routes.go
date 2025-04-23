package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"queue-system/internal/models"
	"queue-system/internal/service"
)

type Handlers struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handlers {
	return &Handlers{service: s}
}

func (h *Handlers) CallNextClient(w http.ResponseWriter, r *http.Request) {
	var req models.Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	queue, err := h.service.CallNextClient(req.WindowId)
	if err != nil {
		// Обрабатываем случай, когда очередь пуста или окно занято
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := models.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Клиент вызван к окну",
		Data: models.Data{
			Id:            queue.ID,
			Number:        queue.Number,
			ServiceTypeId: queue.ServiceTypeID,
			WindowId:      queue.WindowID,
			Status:        queue.Status,
			CreatedAt:     queue.CreatedAt,
			ServedAt:      queue.ServedAt,
			FinishedAt:    queue.FinishedAt,
		},
	}

	writeJSON(w, resp)
}

func (h *Handlers) FinishedClient(w http.ResponseWriter, r *http.Request) {
	var req models.Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	queue, err := h.service.FinishClient(req.WindowId)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := models.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Обслуживание клиента завершено",
		Data: models.Data{
			Id:            queue.ID,
			Number:        queue.Number,
			ServiceTypeId: queue.ServiceTypeID,
			WindowId:      queue.WindowID,
			Status:        queue.Status,
			CreatedAt:     queue.CreatedAt,
			ServedAt:      queue.ServedAt,
			FinishedAt:    queue.FinishedAt,
		},
	}

	writeJSON(w, resp)
}

func (h *Handlers) GetCurrentClient(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	client, err := h.service.GetCurrentClient(req.WindowId)
	if err != nil {
		writeError(w, "Ошибка при получении клиента", http.StatusInternalServerError)
		return
	}

	if client == nil {
		writeError(w, "Нет активного клиента у этого окна", http.StatusNotFound)
		return
	}

	resp := models.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Текущий клиент найден",
		Data: models.Data{
			Id:            client.ID,
			Number:        client.Number,
			ServiceTypeId: client.ServiceTypeID,
			WindowId:      client.WindowID,
			Status:        client.Status,
			CreatedAt:     client.CreatedAt,
			ServedAt:      client.ServedAt,
			FinishedAt:    client.FinishedAt,
		},
	}

	writeJSON(w, resp)
}

func (h *Handlers) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterAndQueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	// Валидация номера
	if len(req.Phone) > 9 || !isDigitsOnly(req.Phone) {
		writeError(w, "Номер телефона должен содержать только цифры и не превышать 9 знаков", http.StatusBadRequest)
		return
	}

	queue, err := h.service.CreateTicketForClient(req.Phone, req.ServiceTypeId)
	if err != nil {
		if errors.Is(err, models.ErrAlreadyInQueue) {
			// Уже в очереди
			resp := models.Response{
				Code:    http.StatusOK,
				Success: true,
				Message: "Клиент уже находится в очереди",
				Data: models.Data{
					Id:            queue.ID,
					Number:        queue.Number,
					ServiceTypeId: queue.ServiceTypeID,
					Phone:         req.Phone,
					WindowId:      queue.WindowID,
					Status:        queue.Status,
					CreatedAt:     queue.CreatedAt,
					ServedAt:      queue.ServedAt,
					FinishedAt:    queue.FinishedAt,
				},
			}
			writeJSON(w, resp)
			return
		}

		writeError(w, "Ошибка при регистрации и постановке в очередь", http.StatusInternalServerError)
		return
	}

	resp := models.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Клиент зарегистрирован и поставлен в очередь",
		Data: models.Data{
			Id:            queue.ID,
			Number:        queue.Number,
			ServiceTypeId: queue.ServiceTypeID,
			Phone:         req.Phone,
			WindowId:      queue.WindowID,
			Status:        queue.Status,
			CreatedAt:     queue.CreatedAt,
			ServedAt:      queue.ServedAt,
			FinishedAt:    queue.FinishedAt,
		},
	}
	writeJSON(w, resp)

}
func (h *Handlers) DashboardState(w http.ResponseWriter, _ *http.Request) {

	calls, err := h.service.ListActiveCalls()
	if err != nil {
		writeError(w, "Не удалось получить состояние", 500)
		return
	}
	type BoardResponse struct {
		Code    int               `json:"code"`
		Success bool              `json:"success"`
		Calls   []models.CallInfo `json:"calls"`
	}
	resp := BoardResponse{200, true, calls}
	json.NewEncoder(w).Encode(resp)
}

func isDigitsOnly(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func writeJSON(w http.ResponseWriter, resp models.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Ошибка при отправке JSON: %v", err)
	}
}

func writeError(w http.ResponseWriter, message string, code int) {
	resp := models.Response{
		Code:    code,
		Success: false,
		Message: message,
	}
	writeJSON(w, resp)
}
