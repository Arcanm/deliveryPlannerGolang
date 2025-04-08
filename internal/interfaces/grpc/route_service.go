package grpc

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Arcanm/deliveryPlannerGolang/internal/application/services"
	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/models"
	"github.com/Arcanm/deliveryPlannerGolang/proto"
)

// RouteService implements the gRPC route service
type RouteService struct {
	proto.UnimplementedRouteServiceServer
	service *services.RouteService
}

// NewRouteService creates a new gRPC route service
func NewRouteService(service *services.RouteService) *RouteService {
	return &RouteService{
		service: service,
	}
}

// CreateRoute creates a new route
func (s *RouteService) CreateRoute(ctx context.Context, req *proto.CreateRouteRequest) (*proto.CreateRouteResponse, error) {
	driverID, err := primitive.ObjectIDFromHex(req.DriverId)
	if err != nil {
		return nil, err
	}

	date := req.Date.AsTime()
	route, err := s.service.CreateRoute(ctx, driverID, date)
	if err != nil {
		return nil, err
	}

	return &proto.CreateRouteResponse{
		Route: convertRouteToProtoResponse(route),
	}, nil
}

// GetRoute retrieves a route by ID
func (s *RouteService) GetRoute(ctx context.Context, req *proto.GetRouteRequest) (*proto.GetRouteResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	route, err := s.service.GetRoute(ctx, id)
	if err != nil {
		return nil, err
	}
	if route == nil {
		return nil, nil
	}

	return &proto.GetRouteResponse{
		Route: convertRouteToProtoResponse(route),
	}, nil
}

// ListRoutes retrieves all routes
func (s *RouteService) ListRoutes(ctx context.Context, req *proto.ListRoutesRequest) (*proto.ListRoutesResponse, error) {
	routes, err := s.service.ListRoutes(ctx)
	if err != nil {
		return nil, err
	}

	protoRoutes := make([]*proto.Route, len(routes))
	for i, route := range routes {
		protoRoutes[i] = convertRouteToProtoResponse(route)
	}

	return &proto.ListRoutesResponse{
		Routes: protoRoutes,
	}, nil
}

// UpdateRoute updates a route
func (s *RouteService) UpdateRoute(ctx context.Context, req *proto.UpdateRouteRequest) (*proto.UpdateRouteResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	route, err := s.service.GetRoute(ctx, id)
	if err != nil {
		return nil, err
	}
	if route == nil {
		return nil, nil
	}

	driverID, err := primitive.ObjectIDFromHex(req.DriverId)
	if err != nil {
		return nil, err
	}

	route.DriverID = driverID
	route.Date = req.Date.AsTime()
	route.EstimatedDistanceKm = float64(req.EstimatedDistanceKm)
	route.EstimatedTimeMin = int(req.EstimatedTimeMin)

	if err := s.service.UpdateRoute(ctx, route); err != nil {
		return nil, err
	}

	return &proto.UpdateRouteResponse{
		Route: convertRouteToProtoResponse(route),
	}, nil
}

// MarkRouteAsCompleted marks a route as completed
func (s *RouteService) MarkRouteAsCompleted(ctx context.Context, req *proto.MarkRouteAsCompletedRequest) (*proto.MarkRouteAsCompletedResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	route, err := s.service.GetRoute(ctx, id)
	if err != nil {
		return nil, err
	}

	route.Status = models.RouteStatusCompleted
	if err := s.service.UpdateRoute(ctx, route); err != nil {
		return nil, err
	}

	return &proto.MarkRouteAsCompletedResponse{
		Route: convertRouteToProtoResponse(route),
	}, nil
}

// AddPackagesToRoute adds packages to a route
func (s *RouteService) AddPackagesToRoute(ctx context.Context, req *proto.AddPackagesToRouteRequest) (*proto.AddPackagesToRouteResponse, error) {
	routeID, err := primitive.ObjectIDFromHex(req.RouteId)
	if err != nil {
		return nil, err
	}

	packageIDs := make([]primitive.ObjectID, len(req.PackageIds))
	for i, id := range req.PackageIds {
		packageID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		packageIDs[i] = packageID
	}

	if err := s.service.AddPackagesToRoute(ctx, routeID, packageIDs); err != nil {
		return nil, err
	}

	return &proto.AddPackagesToRouteResponse{}, nil
}

// UpdatePackageDeliveryStatus updates a package's delivery status in a route
func (s *RouteService) UpdatePackageDeliveryStatus(ctx context.Context, req *proto.UpdatePackageDeliveryStatusRequest) (*proto.UpdatePackageDeliveryStatusResponse, error) {
	routeID, err := primitive.ObjectIDFromHex(req.RouteId)
	if err != nil {
		return nil, err
	}

	packageID, err := primitive.ObjectIDFromHex(req.PackageId)
	if err != nil {
		return nil, err
	}

	if err := s.service.UpdatePackageDeliveryStatus(ctx, routeID, packageID, req.Delivered); err != nil {
		return nil, err
	}

	return &proto.UpdatePackageDeliveryStatusResponse{}, nil
}

// DeleteRoute deletes a route
func (s *RouteService) DeleteRoute(ctx context.Context, req *proto.DeleteRouteRequest) (*proto.DeleteRouteResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	if err := s.service.DeleteRoute(ctx, id); err != nil {
		return nil, err
	}

	return &proto.DeleteRouteResponse{}, nil
}

// Helper functions to convert between domain and proto models
func convertRouteToProtoResponse(route *models.Route) *proto.Route {
	if route == nil {
		return nil
	}

	packages := make([]*proto.PackageRoute, len(route.Packages))
	for i, pkg := range route.Packages {
		packages[i] = &proto.PackageRoute{
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
		Packages:            packages,
		EstimatedDistanceKm: float32(route.EstimatedDistanceKm),
		EstimatedTimeMin:    int32(route.EstimatedTimeMin),
		Completed:           route.Status == models.RouteStatusCompleted,
		CreatedAt:           timestamppb.New(route.CreatedAt),
		UpdatedAt:           timestamppb.New(route.UpdatedAt),
	}
}
