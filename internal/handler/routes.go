package handler

import (
	"github.com/gorilla/mux"
	"queue-system/internal/service"
)

type Handlers struct {
	svc service.IService
}

func NewHandler(s service.IService) *Handlers {
	return &Handlers{svc: s}
}

func (h *Handlers) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/dashboard/state", h.DashboardState).Methods("GET")
	router.HandleFunc("/queue/call", h.CallNextClient).Methods("POST")
	router.HandleFunc("/client/register-and-queue", h.CreateTicket).Methods("POST")
	router.HandleFunc("/queue/finish", h.FinishedClient).Methods("POST")
	router.HandleFunc("/queue/listClient", h.GetCurrentClient).Methods("GET")

	return router
}
