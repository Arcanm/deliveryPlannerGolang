package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/models"
	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/repositories"
)

// PackageService handles package business logic
type PackageService struct {
	packageRepo *repositories.PackageRepository
}

// NewPackageService creates a new package service
func NewPackageService(packageRepo *repositories.PackageRepository) *PackageService {
	return &PackageService{
		packageRepo: packageRepo,
	}
}

// CreatePackage creates a new package
func (s *PackageService) CreatePackage(ctx context.Context, trackingNumber, customerName, customerAddress, customerPhone string, weightKg, volumeM3 float64) (*models.Package, error) {
	pkg := models.NewPackage(trackingNumber, customerName, customerAddress, customerPhone, weightKg, volumeM3)

	if err := s.packageRepo.Create(ctx, pkg); err != nil {
		return nil, err
	}

	return pkg, nil
}

// GetPackage retrieves a package by ID
func (s *PackageService) GetPackage(ctx context.Context, id primitive.ObjectID) (*models.Package, error) {
	return s.packageRepo.GetByID(ctx, id)
}

// GetPackageByTrackingNumber retrieves a package by tracking number
func (s *PackageService) GetPackageByTrackingNumber(ctx context.Context, trackingNumber string) (*models.Package, error) {
	return s.packageRepo.GetByTrackingNumber(ctx, trackingNumber)
}

// ListPackages retrieves all packages
func (s *PackageService) ListPackages(ctx context.Context) ([]*models.Package, error) {
	return s.packageRepo.List(ctx)
}

// UpdatePackage updates a package
func (s *PackageService) UpdatePackage(ctx context.Context, id primitive.ObjectID, trackingNumber, customerName, customerAddress, customerPhone string, weightKg, volumeM3 float64) (*models.Package, error) {
	pkg, err := s.packageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	pkg.TrackingNumber = trackingNumber
	pkg.CustomerName = customerName
	pkg.CustomerAddress = customerAddress
	pkg.CustomerPhone = customerPhone
	pkg.WeightKg = weightKg
	pkg.VolumeM3 = volumeM3
	pkg.UpdatedAt = time.Now()

	if err := s.packageRepo.Update(ctx, pkg); err != nil {
		return nil, err
	}

	return pkg, nil
}

// DeletePackage deletes a package
func (s *PackageService) DeletePackage(ctx context.Context, id primitive.ObjectID) error {
	pkg, err := s.packageRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if pkg.Delivered {
		return errors.New("cannot delete a delivered package")
	}

	return s.packageRepo.Delete(ctx, id)
}

// MarkAsDelivered marks a package as delivered
func (s *PackageService) MarkAsDelivered(ctx context.Context, id primitive.ObjectID) (*models.Package, error) {
	pkg, err := s.packageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	pkg.MarkAsDelivered()

	if err := s.packageRepo.Update(ctx, pkg); err != nil {
		return nil, err
	}

	return pkg, nil
}
