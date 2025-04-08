package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/models"
)

type PackageRepository struct {
	collection *mongo.Collection
}

func NewPackageRepository(db *mongo.Database) *PackageRepository {
	return &PackageRepository{
		collection: db.Collection("packages"),
	}
}

func (r *PackageRepository) Create(ctx context.Context, pkg *models.Package) error {
	pkg.CreatedAt = time.Now()
	pkg.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, pkg)
	if err != nil {
		return err
	}

	pkg.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *PackageRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Package, error) {
	var pkg models.Package
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&pkg)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *PackageRepository) GetByTrackingNumber(ctx context.Context, trackingNumber string) (*models.Package, error) {
	var pkg models.Package
	err := r.collection.FindOne(ctx, bson.M{"tracking_number": trackingNumber}).Decode(&pkg)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *PackageRepository) List(ctx context.Context) ([]*models.Package, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var packages []*models.Package
	if err = cursor.All(ctx, &packages); err != nil {
		return nil, err
	}
	return packages, nil
}

func (r *PackageRepository) Update(ctx context.Context, pkg *models.Package) error {
	pkg.UpdatedAt = time.Now()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": pkg.ID}, pkg)
	return err
}

func (r *PackageRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *PackageRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.PackageStatus) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}
