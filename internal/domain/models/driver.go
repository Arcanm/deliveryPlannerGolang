package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VehicleType represents the type of vehicle a driver uses
type VehicleType string

const (
	VehicleTypeBike  VehicleType = "bike"
	VehicleTypeVan   VehicleType = "van"
	VehicleTypeTruck VehicleType = "truck"
)

// Driver represents a delivery driver
type Driver struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	VehicleType VehicleType        `bson:"vehicle_type"`
	Active      bool               `bson:"active"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

// NewDriver creates a new driver instance
func NewDriver(name string, vehicleType VehicleType) *Driver {
	now := time.Now()
	return &Driver{
		Name:        name,
		VehicleType: vehicleType,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Validate performs basic validation on the driver
func (d *Driver) Validate() error {
	// TODO: Implement validation logic
	return nil
}
