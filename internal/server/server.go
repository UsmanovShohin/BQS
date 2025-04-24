package server

import (
	"context"
	"fmt"
	"net/http"
	"queue-system/internal/config"
	"queue-system/internal/handler"
	"queue-system/pkg"
)

type Server struct {
	server *http.Server
	cfg    *config.Configs
}

func NewServer(cfg *config.Configs, handler *handler.Handlers) *Server {
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port), // гарантированно в формате ":8080"
		Handler:        handler.InitRoutes(),
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	return &Server{
		server: s,
		cfg:    cfg,
	}
}

func (s *Server) Start() error {
	pkg.Log.Infof("HTTP слушает :%d", s.cfg.Server.Port)
	return s.server.ListenAndServe()
}

//func (s *Server) Start() error {
//	go func() {
//		pkg.Log.Infof("Сервер запущен на порту %d", s.cfg.Server.Port)
//		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//			pkg.Log.Fatalf("Ошибка запуска сервера: %v", err)
//		}
//	}()
//
//	// Обработка сигналов остановки
//	quit := make(chan os.Signal, 1)
//	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
//	<-quit
//
//	pkg.Log.Info("Остановка сервера...")
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	if err := s.Shutdown(ctx); err != nil {
//		pkg.Log.Errorf("Ошибка при остановке сервера: %v", err)
//		return err
//	}
//
//	pkg.Log.Info("Сервер остановлен корректно")
//	return nil
//}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
