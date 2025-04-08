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

// DriverService implements the gRPC driver service
type DriverService struct {
	proto.UnimplementedDriverServiceServer
	service *services.DriverService
}

// NewDriverService creates a new gRPC driver service
func NewDriverService(service *services.DriverService) *DriverService {
	return &DriverService{
		service: service,
	}
}

// CreateDriver creates a new driver
func (s *DriverService) CreateDriver(ctx context.Context, req *proto.CreateDriverRequest) (*proto.CreateDriverResponse, error) {
	vehicleType := models.VehicleType(req.VehicleType.String())
	driver, err := s.service.CreateDriver(ctx, req.Name, vehicleType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create driver: %v", err)
	}

	return &proto.CreateDriverResponse{
		Driver: convertDriverToProto(driver),
	}, nil
}

// GetDriver retrieves a driver by ID
func (s *DriverService) GetDriver(ctx context.Context, req *proto.GetDriverRequest) (*proto.GetDriverResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid driver id: %v", err)
	}

	driver, err := s.service.GetDriver(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "driver not found: %v", err)
	}

	return &proto.GetDriverResponse{
		Driver: convertDriverToProto(driver),
	}, nil
}

// ListDrivers retrieves all drivers
func (s *DriverService) ListDrivers(ctx context.Context, req *proto.ListDriversRequest) (*proto.ListDriversResponse, error) {
	drivers, err := s.service.ListDrivers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list drivers: %v", err)
	}

	protoDrivers := make([]*proto.Driver, len(drivers))
	for i, driver := range drivers {
		protoDrivers[i] = convertDriverToProto(driver)
	}

	return &proto.ListDriversResponse{
		Drivers: protoDrivers,
	}, nil
}

// UpdateDriver updates a driver
func (s *DriverService) UpdateDriver(ctx context.Context, req *proto.UpdateDriverRequest) (*proto.UpdateDriverResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid driver id: %v", err)
	}

	vehicleType := models.VehicleType(req.VehicleType.String())
	driver, err := s.service.UpdateDriver(ctx, id, req.Name, vehicleType, req.Active)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update driver: %v", err)
	}

	return &proto.UpdateDriverResponse{
		Driver: convertDriverToProto(driver),
	}, nil
}

// DeleteDriver deletes a driver
func (s *DriverService) DeleteDriver(ctx context.Context, req *proto.DeleteDriverRequest) (*proto.DeleteDriverResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid driver id: %v", err)
	}

	if err := s.service.DeleteDriver(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete driver: %v", err)
	}

	return &proto.DeleteDriverResponse{}, nil
}

// GetDriverRoutes retrieves all routes for a driver
func (s *DriverService) GetDriverRoutes(ctx context.Context, req *proto.GetDriverRoutesRequest) (*proto.GetDriverRoutesResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.DriverId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid driver id: %v", err)
	}

	routes, err := s.service.GetDriverRoutes(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get driver routes: %v", err)
	}

	protoRoutes := make([]*proto.Route, len(routes))
	for i, route := range routes {
		protoRoutes[i] = convertRouteToProto(route)
	}

	return &proto.GetDriverRoutesResponse{
		Routes: protoRoutes,
	}, nil
}

// Helper functions to convert between domain and proto models
func convertDriverToProto(driver *models.Driver) *proto.Driver {
	if driver == nil {
		return nil
	}

	return &proto.Driver{
		Id:          driver.ID.Hex(),
		Name:        driver.Name,
		VehicleType: proto.VehicleType(proto.VehicleType_value[string(driver.VehicleType)]),
		Active:      driver.Active,
		CreatedAt:   timestamppb.New(driver.CreatedAt),
		UpdatedAt:   timestamppb.New(driver.UpdatedAt),
	}
}

func convertRouteToProto(route *models.Route) *proto.Route {
	protoPackages := make([]*proto.PackageRoute, len(route.Packages))
	for i, pkg := range route.Packages {
		protoPackages[i] = &proto.PackageRoute{
			PackageId:         pkg.PackageID.Hex(),
			OrderInRoute:      int32(pkg.OrderInRoute),
			Delivered:         pkg.Delivered,
			DeliveryTimestamp: timestamppb.New(*pkg.DeliveryTimestamp),
		}
	}

	return &proto.Route{
		Id:                  route.ID.Hex(),
		DriverId:            route.DriverID.Hex(),
		Date:                timestamppb.New(route.Date),
		Packages:            protoPackages,
		EstimatedDistanceKm: float32(route.EstimatedDistanceKm),
		EstimatedTimeMin:    int32(route.EstimatedTimeMin),
		Completed:           route.Status == models.RouteStatusCompleted,
		CreatedAt:           timestamppb.New(route.CreatedAt),
		UpdatedAt:           timestamppb.New(route.UpdatedAt),
	}
}
