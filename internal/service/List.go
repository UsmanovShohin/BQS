package service

import "queue-system/internal/models"

func (s *Service) ListActiveCalls() ([]models.CallInfo, error) {
	return s.database.ListActiveCalls()
}
