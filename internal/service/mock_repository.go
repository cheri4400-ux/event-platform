package service

import "event-platform/internal/models"

type MockRepo struct{}

func (m *MockRepo) Save(
	event models.Event,
) error {
	return nil
}

func (m *MockRepo) GetAll() (
	[]models.Event,
	error,
) {
	return nil, nil
}
