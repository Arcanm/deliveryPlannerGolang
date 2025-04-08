package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PackageStatus represents the current status of a package
type PackageStatus string

const (
	PackageStatusPending   PackageStatus = "pending"
	PackageStatusAssigned  PackageStatus = "assigned"
	PackageStatusDelivered PackageStatus = "delivered"
	PackageStatusCancelled PackageStatus = "cancelled"
)

// Location represents a geographical location
type Location struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
}

// Package represents a delivery package
type Package struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TrackingNumber    string             `bson:"tracking_number" json:"tracking_number"`
	CustomerName      string             `bson:"customer_name" json:"customer_name"`
	CustomerAddress   string             `bson:"customer_address" json:"customer_address"`
	CustomerPhone     string             `bson:"customer_phone" json:"customer_phone"`
	WeightKg          float64            `bson:"weight_kg" json:"weight_kg"`
	VolumeM3          float64            `bson:"volume_m3" json:"volume_m3"`
	Delivered         bool               `bson:"delivered" json:"delivered"`
	DeliveryTimestamp *time.Time         `bson:"delivery_timestamp,omitempty" json:"delivery_timestamp,omitempty"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewPackage creates a new package instance
func NewPackage(trackingNumber, customerName, customerAddress, customerPhone string, weightKg, volumeM3 float64) *Package {
	now := time.Now()
	return &Package{
		TrackingNumber:    trackingNumber,
		CustomerName:      customerName,
		CustomerAddress:   customerAddress,
		CustomerPhone:     customerPhone,
		WeightKg:          weightKg,
		VolumeM3:          volumeM3,
		Delivered:         false,
		DeliveryTimestamp: nil,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// MarkAsDelivered marks the package as delivered
func (p *Package) MarkAsDelivered() {
	now := time.Now()
	p.Delivered = true
	p.DeliveryTimestamp = &now
	p.UpdatedAt = now
}

// Validate performs basic validation on the package
func (p *Package) Validate() error {
	// TODO: Implement validation logic
	return nil
}
