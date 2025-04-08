package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/models"
)

type DriverRepository struct {
	collection *mongo.Collection
}

func NewDriverRepository(db *mongo.Database) *DriverRepository {
	return &DriverRepository{
		collection: db.Collection("drivers"),
	}
}

func (r *DriverRepository) Create(ctx context.Context, driver *models.Driver) error {
	driver.CreatedAt = time.Now()
	driver.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, driver)
	if err != nil {
		return err
	}

	driver.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *DriverRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Driver, error) {
	var driver models.Driver
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&driver)
	if err != nil {
		return nil, err
	}
	return &driver, nil
}

func (r *DriverRepository) List(ctx context.Context) ([]*models.Driver, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var drivers []*models.Driver
	if err = cursor.All(ctx, &drivers); err != nil {
		return nil, err
	}
	return drivers, nil
}

func (r *DriverRepository) Update(ctx context.Context, driver *models.Driver) error {
	driver.UpdatedAt = time.Now()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": driver.ID}, driver)
	return err
}

func (r *DriverRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
