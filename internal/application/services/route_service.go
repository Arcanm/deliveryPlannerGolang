package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/models"
	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/repositories"
)

// RouteService handles route business logic
type RouteService struct {
	routeRepo   *repositories.RouteRepository
	driverRepo  *repositories.DriverRepository
	packageRepo *repositories.PackageRepository
}

// NewRouteService creates a new route service
func NewRouteService(routeRepo *repositories.RouteRepository, driverRepo *repositories.DriverRepository, packageRepo *repositories.PackageRepository) *RouteService {
	return &RouteService{
		routeRepo:   routeRepo,
		driverRepo:  driverRepo,
		packageRepo: packageRepo,
	}
}

// CreateRoute creates a new route for a driver
func (s *RouteService) CreateRoute(ctx context.Context, driverID primitive.ObjectID, date time.Time) (*models.Route, error) {
	// Verify driver exists and is active
	driver, err := s.driverRepo.GetByID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, fmt.Errorf("driver not found")
	}
	if !driver.Active {
		return nil, fmt.Errorf("driver is not active")
	}

	route := models.NewRoute(driverID, date)
	if err := route.Validate(); err != nil {
		return nil, err
	}

	if err := s.routeRepo.Create(ctx, route); err != nil {
		return nil, err
	}

	return route, nil
}

// GetRoute retrieves a route by ID
func (s *RouteService) GetRoute(ctx context.Context, id primitive.ObjectID) (*models.Route, error) {
	return s.routeRepo.GetByID(ctx, id)
}

// GetDriverRoutes retrieves all routes for a driver
func (s *RouteService) GetDriverRoutes(ctx context.Context, driverID primitive.ObjectID) ([]*models.Route, error) {
	return s.routeRepo.GetByDriverID(ctx, driverID)
}

// ListRoutes retrieves all routes
func (s *RouteService) ListRoutes(ctx context.Context) ([]*models.Route, error) {
	return s.routeRepo.List(ctx)
}

// AddPackagesToRoute adds packages to a route and calculates the route
func (s *RouteService) AddPackagesToRoute(ctx context.Context, routeID primitive.ObjectID, packageIDs []primitive.ObjectID) error {
	// Get the route
	route, err := s.routeRepo.GetByID(ctx, routeID)
	if err != nil {
		return err
	}
	if route == nil {
		return fmt.Errorf("route not found")
	}

	// Verify route is in planned status
	if route.Status != models.RouteStatusPending {
		return fmt.Errorf("can only add packages to pending routes")
	}

	// Get all packages
	var packages []*models.Package
	for _, id := range packageIDs {
		pkg, err := s.packageRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if pkg == nil {
			return fmt.Errorf("package %s not found", id.Hex())
		}
		if pkg.Delivered {
			return fmt.Errorf("package %s is already delivered", id.Hex())
		}
		packages = append(packages, pkg)
	}

	// Calculate optimal route (for now, just add packages in order)
	// TODO: Implement more sophisticated route optimization
	for _, pkg := range packages {
		route.AddPackage(pkg.ID)
	}

	// Calculate estimated distance and time
	// TODO: Implement actual distance calculation
	route.EstimatedDistanceKm = calculateEstimatedDistance(packages)
	route.EstimatedTimeMin = calculateEstimatedTime(route.EstimatedDistanceKm)

	// Update the route
	return s.routeRepo.Update(ctx, route)
}

// UpdateRouteStatus updates a route's status
func (s *RouteService) UpdateRouteStatus(ctx context.Context, id primitive.ObjectID, status models.RouteStatus) error {
	route, err := s.routeRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if route == nil {
		return fmt.Errorf("route not found")
	}

	// Validate status transition
	if route.Status == models.RouteStatusCompleted && status != models.RouteStatusCompleted {
		return fmt.Errorf("cannot change status of a completed route")
	}

	// Update route status
	if err := route.UpdateStatus(status); err != nil {
		return err
	}

	// Update package statuses based on route status
	if status == models.RouteStatusActive {
		for _, pkg := range route.Packages {
			if err := s.packageRepo.UpdateStatus(ctx, pkg.PackageID, models.PackageStatusAssigned); err != nil {
				return err
			}
		}
	} else if status == models.RouteStatusCompleted {
		for _, pkg := range route.Packages {
			if err := s.packageRepo.UpdateStatus(ctx, pkg.PackageID, models.PackageStatusDelivered); err != nil {
				return err
			}
		}
	}

	return s.routeRepo.UpdateStatus(ctx, id, status)
}

// UpdatePackageDeliveryStatus updates a package's delivery status in a route
func (s *RouteService) UpdatePackageDeliveryStatus(ctx context.Context, routeID, packageID primitive.ObjectID, delivered bool) error {
	route, err := s.routeRepo.GetByID(ctx, routeID)
	if err != nil {
		return err
	}
	if route == nil {
		return fmt.Errorf("route not found")
	}

	// Verify route is in progress
	if route.Status != models.RouteStatusActive {
		return fmt.Errorf("can only update package status for routes in progress")
	}

	// Update package status in route
	if err := s.routeRepo.UpdatePackageStatus(ctx, routeID, packageID, delivered); err != nil {
		return err
	}

	// Update package status
	status := models.PackageStatusAssigned
	if delivered {
		status = models.PackageStatusDelivered
	}
	return s.packageRepo.UpdateStatus(ctx, packageID, status)
}

// Helper functions for route calculations
func calculateEstimatedDistance(packages []*models.Package) float64 {
	// TODO: Implement actual distance calculation using coordinates
	// For now, return a mock value
	return float64(len(packages)) * 5.0 // 5km per package
}

func calculateEstimatedTime(distanceKm float64) int {
	// TODO: Implement actual time calculation
	// For now, assume average speed of 50 km/h
	return int(distanceKm / 50.0 * 60.0) // Convert to minutes
}

// UpdateRoute updates an existing route
func (s *RouteService) UpdateRoute(ctx context.Context, route *models.Route) error {
	return s.routeRepo.Update(ctx, route)
}

// DeleteRoute deletes a route by ID
func (s *RouteService) DeleteRoute(ctx context.Context, id primitive.ObjectID) error {
	return s.routeRepo.Delete(ctx, id)
}
