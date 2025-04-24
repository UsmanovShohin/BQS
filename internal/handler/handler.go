package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"queue-system/internal/models"
	"queue-system/pkg"
)

func (h *Handlers) CallNextClient(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	queue, err := h.svc.CallNextClient(req.WindowId)
	if err != nil {
		pkg.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	pkg.WriteSuccess(w, "Клиент вызван к окну", queue)
}

func (h *Handlers) FinishedClient(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	queue, err := h.svc.FinishClient(req.WindowId)
	if err != nil {
		pkg.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	pkg.WriteSuccess(w, "Обслуживание клиента завершено", queue)
}

func (h *Handlers) GetCurrentClient(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	client, err := h.svc.GetCurrentClient(req.WindowId)
	if err != nil {
		pkg.WriteError(w, "Ошибка при получении клиента", http.StatusInternalServerError)
		return
	}
	if client == nil {
		pkg.WriteError(w, "Нет активного клиента у этого окна", http.StatusNotFound)
		return
	}

	pkg.WriteSuccess(w, "Текущий клиент найден", client)
}

func (h *Handlers) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterAndQueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	if err := pkg.ValidatePhone(req.Phone); err != nil {
		pkg.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	queue, err := h.svc.CreateTicketForClient(req.Phone, req.ServiceTypeId)
	if err != nil {
		if errors.Is(err, models.ErrAlreadyInQueue) {
			pkg.WriteSuccess(w, "Клиент уже находится в очереди", queue)
			return
		}
		pkg.WriteError(w, "Ошибка при регистрации и постановке в очередь", http.StatusInternalServerError)
		return
	}

	pkg.WriteSuccess(w, "Клиент зарегистрирован и поставлен в очередь", queue)
}

func (h *Handlers) DashboardState(w http.ResponseWriter, _ *http.Request) {
	calls, err := h.svc.ListActiveCalls()
	if err != nil {
		pkg.WriteError(w, "Не удалось получить состояние", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(struct {
		Code    int               `json:"code"`
		Success bool              `json:"success"`
		Calls   []models.CallInfo `json:"calls"`
	}{
		Code:    http.StatusOK,
		Success: true,
		Calls:   calls,
	})
}
