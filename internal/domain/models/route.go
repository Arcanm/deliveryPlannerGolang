package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RouteStatus represents the current status of a route
type RouteStatus string

const (
	RouteStatusPending   RouteStatus = "pending"
	RouteStatusActive    RouteStatus = "active"
	RouteStatusCompleted RouteStatus = "completed"
	RouteStatusCancelled RouteStatus = "cancelled"
)

// PackageRoute represents a package assigned to a route
type PackageRoute struct {
	PackageID         primitive.ObjectID `bson:"package_id"`
	OrderInRoute      int                `bson:"order_in_route"`
	Delivered         bool               `bson:"delivered"`
	DeliveryTimestamp *time.Time         `bson:"delivery_timestamp,omitempty"`
}

// Route represents a delivery route
type Route struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	DriverID            primitive.ObjectID `bson:"driver_id"`
	Date                time.Time          `bson:"date"`
	Packages            []PackageRoute     `bson:"packages"`
	EstimatedDistanceKm float64            `bson:"estimated_distance_km"`
	EstimatedTimeMin    int                `bson:"estimated_time_min"`
	Status              RouteStatus        `bson:"status"`
	CreatedAt           time.Time          `bson:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at"`
}

// NewRoute creates a new route instance
func NewRoute(driverID primitive.ObjectID, date time.Time) *Route {
	now := time.Now()
	return &Route{
		DriverID:            driverID,
		Date:                date,
		Packages:            make([]PackageRoute, 0),
		EstimatedDistanceKm: 0,
		EstimatedTimeMin:    0,
		Status:              RouteStatusPending,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}

// AddPackage adds a package to the route
func (r *Route) AddPackage(packageID primitive.ObjectID) {
	r.Packages = append(r.Packages, PackageRoute{
		PackageID:    packageID,
		OrderInRoute: len(r.Packages) + 1,
		Delivered:    false,
	})
	r.UpdatedAt = time.Now()
}

// UpdatePackageStatus updates the delivery status of a package in the route
func (r *Route) UpdatePackageStatus(packageID primitive.ObjectID, delivered bool) bool {
	for i := range r.Packages {
		if r.Packages[i].PackageID == packageID {
			r.Packages[i].Delivered = delivered
			if delivered {
				now := time.Now()
				r.Packages[i].DeliveryTimestamp = &now
			} else {
				r.Packages[i].DeliveryTimestamp = nil
			}
			r.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// UpdateStatus updates the route status
func (r *Route) UpdateStatus(status RouteStatus) error {
	// Validate that all packages are delivered before completing the route
	if status == RouteStatusCompleted {
		for _, p := range r.Packages {
			if !p.Delivered {
				return ErrRouteHasPendingPackages
			}
		}
	}

	r.Status = status
	r.UpdatedAt = time.Now()
	return nil
}

// Validate performs basic validation on the route
func (r *Route) Validate() error {
	// TODO: Implement validation logic
	return nil
}
