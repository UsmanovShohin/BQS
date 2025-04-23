package service

import (
	"fmt"
	"strconv"
	"strings"
)

func (s *Service) GeneraQueueNumber(serviceTypeId int) (string, error) {
	serviceType, err := s.database.GetServiceTypeById(serviceTypeId)
	if err != nil {
		return "", fmt.Errorf("error getting service type by id: %v", err)
	}
	lastQueue, err := s.database.GetLastByServiceType(serviceTypeId)
	if err != nil {
		return "", fmt.Errorf("error getting last queue by id: %v", err)
	}
	var number int
	if lastQueue == nil {
		number = 1
	} else {
		numStr := strings.TrimPrefix(lastQueue.Number, serviceType.Code)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return "", fmt.Errorf("не правильный фортма номера: %v", err)
		}
		number = num + 1
	}

	queueNumber := fmt.Sprintf("%s%03d", serviceType.Code, number)
	return queueNumber, nil
}
