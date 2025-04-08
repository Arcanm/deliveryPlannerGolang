package grpc

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Arcanm/deliveryPlannerGolang/internal/application/services"
	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/models"
	"github.com/Arcanm/deliveryPlannerGolang/proto"
)

// PackageService implements the gRPC package service
type PackageService struct {
	proto.UnimplementedPackageServiceServer
	service      *services.PackageService
	routeService *services.RouteService
}

// NewPackageService creates a new gRPC package service
func NewPackageService(service *services.PackageService, routeService *services.RouteService) *PackageService {
	return &PackageService{
		service:      service,
		routeService: routeService,
	}
}

// CreatePackage creates a new package
func (s *PackageService) CreatePackage(ctx context.Context, req *proto.CreatePackageRequest) (*proto.CreatePackageResponse, error) {
	pkg, err := s.service.CreatePackage(ctx, req.TrackingNumber, req.CustomerName, req.CustomerAddress, req.CustomerPhone, float64(req.WeightKg), float64(req.VolumeM3))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create package: %v", err)
	}

	return &proto.CreatePackageResponse{
		Package: convertPackageToProto(pkg),
	}, nil
}

// GetPackage retrieves a package by ID
func (s *PackageService) GetPackage(ctx context.Context, req *proto.GetPackageRequest) (*proto.GetPackageResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid package id: %v", err)
	}

	pkg, err := s.service.GetPackage(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "package not found: %v", err)
	}

	return &proto.GetPackageResponse{
		Package: convertPackageToProto(pkg),
	}, nil
}

// GetPackageByTrackingNumber retrieves a package by tracking number
func (s *PackageService) GetPackageByTrackingNumber(ctx context.Context, req *proto.GetPackageByTrackingNumberRequest) (*proto.GetPackageByTrackingNumberResponse, error) {
	pkg, err := s.service.GetPackageByTrackingNumber(ctx, req.TrackingNumber)
	if err != nil {
		return nil, err
	}
	if pkg == nil {
		return nil, nil
	}

	return &proto.GetPackageByTrackingNumberResponse{
		Package: convertPackageToProto(pkg),
	}, nil
}

// ListPackages retrieves all packages
func (s *PackageService) ListPackages(ctx context.Context, req *proto.ListPackagesRequest) (*proto.ListPackagesResponse, error) {
	packages, err := s.service.ListPackages(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list packages: %v", err)
	}

	protoPackages := make([]*proto.Package, len(packages))
	for i, pkg := range packages {
		protoPackages[i] = convertPackageToProto(pkg)
	}

	return &proto.ListPackagesResponse{
		Packages: protoPackages,
	}, nil
}

// UpdatePackage updates a package
func (s *PackageService) UpdatePackage(ctx context.Context, req *proto.UpdatePackageRequest) (*proto.UpdatePackageResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid package id: %v", err)
	}

	pkg, err := s.service.UpdatePackage(ctx, id, req.TrackingNumber, req.CustomerName, req.CustomerAddress, req.CustomerPhone, float64(req.WeightKg), float64(req.VolumeM3))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update package: %v", err)
	}

	return &proto.UpdatePackageResponse{
		Package: convertPackageToProto(pkg),
	}, nil
}

// UpdatePackageStatus updates a package's status
func (s *PackageService) UpdatePackageStatus(ctx context.Context, req *proto.UpdatePackageStatusRequest) (*proto.UpdatePackageStatusResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	pkg, err := s.service.MarkAsDelivered(ctx, id)
	if err != nil {
		return nil, err
	}

	return &proto.UpdatePackageStatusResponse{
		Package: convertPackageToProto(pkg),
	}, nil
}

// DeletePackage deletes a package
func (s *PackageService) DeletePackage(ctx context.Context, req *proto.DeletePackageRequest) (*proto.DeletePackageResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid package id: %v", err)
	}

	if err := s.service.DeletePackage(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete package: %v", err)
	}

	return &proto.DeletePackageResponse{}, nil
}

// AssignToRoute assigns a package to a route
func (s *PackageService) AssignToRoute(ctx context.Context, req *proto.AssignToRouteRequest) (*proto.AssignToRouteResponse, error) {
	packageID, err := primitive.ObjectIDFromHex(req.PackageId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid package id: %v", err)
	}

	routeID, err := primitive.ObjectIDFromHex(req.RouteId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid route id: %v", err)
	}

	if err := s.routeService.AddPackagesToRoute(ctx, routeID, []primitive.ObjectID{packageID}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to assign package to route: %v", err)
	}

	return &proto.AssignToRouteResponse{}, nil
}

// MarkPackageAsDelivered marks a package as delivered
func (s *PackageService) MarkPackageAsDelivered(ctx context.Context, req *proto.MarkPackageAsDeliveredRequest) (*proto.MarkPackageAsDeliveredResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid package id: %v", err)
	}

	pkg, err := s.service.MarkAsDelivered(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mark package as delivered: %v", err)
	}

	return &proto.MarkPackageAsDeliveredResponse{
		Package: convertPackageToProto(pkg),
	}, nil
}

// GetPackagesByRoute retrieves packages by route
func (s *PackageService) GetPackagesByRoute(ctx context.Context, req *proto.GetPackagesByRouteRequest) (*proto.GetPackagesByRouteResponse, error) {
	routeID, err := primitive.ObjectIDFromHex(req.RouteId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid route id: %v", err)
	}

	route, err := s.routeService.GetRoute(ctx, routeID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	packages := make([]*models.Package, len(route.Packages))
	for i, pkg := range route.Packages {
		packageID, err := primitive.ObjectIDFromHex(pkg.PackageID.Hex())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "invalid package id in route: %v", err)
		}
		pkg, err := s.service.GetPackage(ctx, packageID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get package: %v", err)
		}
		packages[i] = pkg
	}

	protoPackages := make([]*proto.Package, len(packages))
	for i, pkg := range packages {
		protoPackages[i] = convertPackageToProto(pkg)
	}

	return &proto.GetPackagesByRouteResponse{
		Packages: protoPackages,
	}, nil
}

// Helper functions to convert between domain and proto models
func convertPackageToProto(pkg *models.Package) *proto.Package {
	if pkg == nil {
		return nil
	}

	var deliveryTimestamp *timestamppb.Timestamp
	if pkg.DeliveryTimestamp != nil {
		deliveryTimestamp = timestamppb.New(*pkg.DeliveryTimestamp)
	}

	return &proto.Package{
		Id:                pkg.ID.Hex(),
		TrackingNumber:    pkg.TrackingNumber,
		CustomerName:      pkg.CustomerName,
		CustomerAddress:   pkg.CustomerAddress,
		CustomerPhone:     pkg.CustomerPhone,
		WeightKg:          float64(pkg.WeightKg),
		VolumeM3:          float64(pkg.VolumeM3),
		Delivered:         pkg.Delivered,
		DeliveryTimestamp: deliveryTimestamp,
		CreatedAt:         timestamppb.New(pkg.CreatedAt),
		UpdatedAt:         timestamppb.New(pkg.UpdatedAt),
	}
}
