package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client *Client
}

func NewRepository(client *Client) *Repository {
	return &Repository{client: client}
}

type JobResultRepository struct {
	*Repository
	collection *mongo.Collection
}

func NewJobResultRepository(repo *Repository) *JobResultRepository {
	return &JobResultRepository{
		Repository: repo,
		collection: repo.client.GetCollection("job_results"),
	}
}

func (r *JobResultRepository) SaveJobResults(userID primitive.ObjectID, queryTitle string, jobIDs []primitive.ObjectID) (*JobResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobResult := &JobResult{
		UserID:     userID,
		Jobs:       jobIDs,
		QueryTitle: queryTitle,
		CreatedAt:  time.Now(),
	}

	result, err := r.collection.InsertOne(ctx, jobResult)
	if err != nil {
		return nil, fmt.Errorf("failed to save job results: %w", err)
	}

	jobResult.ID = result.InsertedID.(primitive.ObjectID)
	return jobResult, nil
}

func (r *JobResultRepository) GetLatestJobResults(userID primitive.ObjectID) (*JobResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	opts := options.FindOne().SetSort(bson.D{{"created_at", -1}})

	var jobResult JobResult
	err := r.collection.FindOne(ctx, filter, opts).Decode(&jobResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no job results found for user")
		}
		return nil, fmt.Errorf("failed to get latest job results: %w", err)
	}

	return &jobResult, nil
}

func (r *JobResultRepository) GetAllJobResults(userID primitive.ObjectID) ([]JobResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get job results: %w", err)
	}
	defer cursor.Close(ctx)

	var results []JobResult
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode job results: %w", err)
	}

	return results, nil
}

type JobRepository struct {
	*Repository
	collection *mongo.Collection
}

func NewJobRepository(repo *Repository) *JobRepository {
	return &JobRepository{
		Repository: repo,
		collection: repo.client.GetCollection("jobs"),
	}
}

func (r *JobRepository) SaveJobs(jobs []Job) ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// convert to interface slice for bulk insert
	docs := make([]interface{}, len(jobs))
	for i, job := range jobs {
		docs[i] = job
	}

	result, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return nil, fmt.Errorf("failed to save jobs: %w", err)
	}

	// convert inserted IDs to ObjectIDs
	ids := make([]primitive.ObjectID, len(result.InsertedIDs))
	for i, id := range result.InsertedIDs {
		ids[i] = id.(primitive.ObjectID)
	}

	return ids, nil
}

func (r *JobRepository) GetJobsByIDs(ids []primitive.ObjectID) ([]Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by IDs: %w", err)
	}
	defer cursor.Close(ctx)

	var jobs []Job
	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, fmt.Errorf("failed to decode jobs: %w", err)
	}

	return jobs, nil
}

type BusinessRepository struct {
	*Repository
	collection *mongo.Collection
}

func NewBusinessRepository(repo *Repository) *BusinessRepository {
	return &BusinessRepository{
		Repository: repo,
		collection: repo.client.GetCollection("businesses"),
	}
}

func (r *BusinessRepository) SaveBusinesses(businesses []Business) ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	docs := make([]interface{}, len(businesses))
	for i, business := range businesses {
		docs[i] = business
	}

	result, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return nil, fmt.Errorf("failed to save businesses: %w", err)
	}

	ids := make([]primitive.ObjectID, len(result.InsertedIDs))
	for i, id := range result.InsertedIDs {
		ids[i] = id.(primitive.ObjectID)
	}

	return ids, nil
}

func (r *BusinessRepository) GetBusinessesByIDs(ids []primitive.ObjectID) ([]Business, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get businesses by IDs: %w", err)
	}
	defer cursor.Close(ctx)

	var businesses []Business
	if err = cursor.All(ctx, &businesses); err != nil {
		return nil, fmt.Errorf("failed to decode businesses: %w", err)
	}

	return businesses, nil
}

type GeoResultRepository struct {
	*Repository
	collection *mongo.Collection
}

func NewGeoResultRepository(repo *Repository) *GeoResultRepository {
	return &GeoResultRepository{
		Repository: repo,
		collection: repo.client.GetCollection("geo_results"),
	}
}

func (r *GeoResultRepository) SaveGeoResult(userID primitive.ObjectID, zip string, radius int) (*GeoResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	geoResult := &GeoResult{
		UserID:    userID,
		Zip:       zip,
		Radius:    radius,
		CreatedAt: time.Now(),
	}

	result, err := r.collection.InsertOne(ctx, geoResult)
	if err != nil {
		return nil, fmt.Errorf("failed to save geo result: %w", err)
	}

	geoResult.ID = result.InsertedID.(primitive.ObjectID)
	return geoResult, nil
}
