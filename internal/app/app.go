package app

import (
	"log"
	"queue-system/internal/config"
	"queue-system/internal/database"
	"queue-system/internal/handler"
	"queue-system/internal/repository"
	"queue-system/internal/server"
	"queue-system/internal/service"
	"queue-system/pkg"
)

type App struct {
	cfg    *config.Configs
	server *server.Server
}

func New() *App {
	pkg.InitLogger()
	pkg.Log.Info("Сервер запускается...")

	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфига: %v", err)
	}

	conn := database.ConnectDB(cfg)
	repo := repository.NewRepository(conn)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)
	srv := server.NewServer(cfg, h)

	return &App{
		cfg:    cfg,
		server: srv,
	}
}

func (a *App) Run() {
	if err := a.server.Start(); err != nil {
		pkg.Log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
