package pkg

import (
	"encoding/json"
	"errors"
	"net/http"
	"queue-system/internal/models"
	"regexp"
)

func ValidatePhone(phone string) error {
	if len(phone) > 9 {
		return errors.New("номер телефона не должен превышать 9 цифр")
	}
	matched, err := regexp.MatchString(`^\d{1,9}$`, phone)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("номер телефона должен содержать только цифры")
	}
	return nil
}

func WriteSuccess(w http.ResponseWriter, message string, ticket *models.QueueTicket) {
	resp := models.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: message,
		Data: models.Data{
			Id:            ticket.ID,
			Number:        ticket.Number,
			ServiceTypeId: ticket.ServiceTypeID,
			WindowId:      ticket.WindowID,
			Status:        ticket.Status,
			CreatedAt:     ticket.CreatedAt,
			ServedAt:      ticket.ServedAt,
			FinishedAt:    ticket.FinishedAt,
		},
	}
	writeJSON(w, resp)
}

func WriteError(w http.ResponseWriter, message string, code int) {
	resp := models.Response{
		Code:    code,
		Success: false,
		Message: message,
	}
	writeJSON(w, resp)
}

func writeJSON(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
	}
}
