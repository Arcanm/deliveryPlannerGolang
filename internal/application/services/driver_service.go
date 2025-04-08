package services

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/models"
	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/repositories"
)

// DriverService handles driver business logic
type DriverService struct {
	driverRepo *repositories.DriverRepository
	routeRepo  *repositories.RouteRepository
}

// NewDriverService creates a new driver service
func NewDriverService(driverRepo *repositories.DriverRepository, routeRepo *repositories.RouteRepository) *DriverService {
	return &DriverService{
		driverRepo: driverRepo,
		routeRepo:  routeRepo,
	}
}

// CreateDriver creates a new driver
func (s *DriverService) CreateDriver(ctx context.Context, name string, vehicleType models.VehicleType) (*models.Driver, error) {
	driver := &models.Driver{
		Name:        name,
		VehicleType: vehicleType,
		Active:      true,
	}

	if err := s.driverRepo.Create(ctx, driver); err != nil {
		return nil, err
	}

	return driver, nil
}

// GetDriver retrieves a driver by ID
func (s *DriverService) GetDriver(ctx context.Context, id primitive.ObjectID) (*models.Driver, error) {
	return s.driverRepo.GetByID(ctx, id)
}

// ListDrivers retrieves all drivers
func (s *DriverService) ListDrivers(ctx context.Context) ([]*models.Driver, error) {
	return s.driverRepo.List(ctx)
}

// UpdateDriver updates a driver
func (s *DriverService) UpdateDriver(ctx context.Context, id primitive.ObjectID, name string, vehicleType models.VehicleType, active bool) (*models.Driver, error) {
	driver, err := s.driverRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	driver.Name = name
	driver.VehicleType = vehicleType
	driver.Active = active

	if err := s.driverRepo.Update(ctx, driver); err != nil {
		return nil, err
	}

	return driver, nil
}

// DeleteDriver deletes a driver
func (s *DriverService) DeleteDriver(ctx context.Context, id primitive.ObjectID) error {
	// Check if driver has any active routes
	routes, err := s.routeRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, route := range routes {
		if route.DriverID == id && route.Status == models.RouteStatusActive {
			return errors.New("cannot delete driver with active routes")
		}
	}

	return s.driverRepo.Delete(ctx, id)
}

// GetDriverRoutes retrieves all routes for a driver
func (s *DriverService) GetDriverRoutes(ctx context.Context, driverID primitive.ObjectID) ([]*models.Route, error) {
	routes, err := s.routeRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	var driverRoutes []*models.Route
	for _, route := range routes {
		if route.DriverID == driverID {
			driverRoutes = append(driverRoutes, route)
		}
	}

	return driverRoutes, nil
}
