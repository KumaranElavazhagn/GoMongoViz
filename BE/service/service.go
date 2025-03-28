package service

import (
	"gomongoviz/model"
	"gomongoviz/repository"
)

// Service implements the business logic layer of the application
// It sits between the handler and repository layers, processing data from the repository
// before passing it to the handlers
type Service struct {
	repo repository.Repository // Repository interface for data access
}

// GetPorts retrieves all ports associated with a specific object ID
// It delegates to the repository layer and returns the port information
func (s *Service) GetPorts(objectID int) (any, error) {
	ports, err := s.repo.GetPorts(objectID)
	if err != nil {
		return nil, err
	}
	return ports, nil
}

// GetUniqueObjectIDs retrieves all unique object IDs from the repository
// This is used to populate selection dropdowns and filters in the UI
func (s *Service) GetUniqueObjectIDs() ([]model.ObjectInfo, error) {
	objectIDs, err := s.repo.GetUniqueObjectIDs()
	if err != nil {
		return nil, err
	}
	return objectIDs, nil
}

// GetDataByObjectID retrieves sensor data for a specific object ID
// Optionally filtered by port number if provided
func (s *Service) GetDataByObjectID(objectID string, portNum string) (*model.SensorDataRes, error) {
	data, err := s.repo.GetDataByObjectID(objectID, portNum)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SaveSensorData saves a batch of sensor data records to the database
// This is used for the CSV upload feature
func (s *Service) SaveSensorData(data []model.SensorData) error {
	return s.repo.SaveSensorData(data)
}

// NewService creates a new service instance with the provided repository
// This follows the dependency injection pattern, allowing for easier testing
func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}
