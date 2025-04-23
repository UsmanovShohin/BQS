package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"queue-system/internal/config"
	"queue-system/internal/database"
	"queue-system/internal/handler"
	"queue-system/internal/service"
	"queue-system/pkg"
)

func main() {
	pkg.InitLogger()
	pkg.Log.Info("Сервер запускается...")
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка при загрузки конфига: %v", err)
	}
	connection := database.ConnectDB(cfg)
	db := database.NewDatabase(connection)
	service1 := service.NewService(db)
	handler1 := handler.NewHandler(service1)

	r := mux.NewRouter()
	r.HandleFunc("/dashboard/state", handler1.DashboardState).Methods("GET")
	r.HandleFunc("/queue/call", handler1.CallNextClient).Methods("POST")
	r.HandleFunc("/client/register-and-queue", handler1.CreateTicket).Methods("POST")
	r.HandleFunc("/queue/finish", handler1.FinishedClient).Methods("POST")
	r.HandleFunc("/queue/listClient", handler1.GetCurrentClient).Methods("GET")
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
