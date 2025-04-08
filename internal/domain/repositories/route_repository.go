package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/models"
)

type RouteRepository struct {
	collection *mongo.Collection
}

func NewRouteRepository(db *mongo.Database) *RouteRepository {
	return &RouteRepository{
		collection: db.Collection("routes"),
	}
}

func (r *RouteRepository) Create(ctx context.Context, route *models.Route) error {
	route.CreatedAt = time.Now()
	route.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, route)
	if err != nil {
		return err
	}

	route.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *RouteRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Route, error) {
	var route models.Route
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&route)
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *RouteRepository) List(ctx context.Context) ([]*models.Route, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var routes []*models.Route
	if err = cursor.All(ctx, &routes); err != nil {
		return nil, err
	}
	return routes, nil
}

func (r *RouteRepository) Update(ctx context.Context, route *models.Route) error {
	route.UpdatedAt = time.Now()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": route.ID}, route)
	return err
}

func (r *RouteRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *RouteRepository) GetByDriverID(ctx context.Context, driverID primitive.ObjectID) ([]*models.Route, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"driver_id": driverID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var routes []*models.Route
	if err = cursor.All(ctx, &routes); err != nil {
		return nil, err
	}
	return routes, nil
}

func (r *RouteRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.RouteStatus) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (r *RouteRepository) UpdatePackageStatus(ctx context.Context, routeID, packageID primitive.ObjectID, delivered bool) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id":          routeID,
			"packages._id": packageID,
		},
		bson.M{
			"$set": bson.M{
				"packages.$.delivered":          delivered,
				"packages.$.delivery_timestamp": time.Now(),
			},
		},
	)
	return err
}
