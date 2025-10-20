package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"password_hash"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type GeoResult struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Zip       string             `bson:"zip" json:"zip"`
	Radius    int                `bson:"radius" json:"radius"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Business struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GeoResultID  primitive.ObjectID `bson:"geo_result_id" json:"geo_result_id"`
	Name         string             `bson:"name" json:"name"`
	Address      string             `bson:"address" json:"address"`
	URL          string             `bson:"url" json:"url"`
	Lat          float64            `bson:"lat" json:"lat"`
	Lon          float64            `bson:"lon" json:"lon"`
}

type Job struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BusinessID  primitive.ObjectID `bson:"business_id" json:"business_id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	URL         string             `bson:"url" json:"url"`
	PostedAt    *time.Time         `bson:"posted_at,omitempty" json:"posted_at,omitempty"`
}

type JobResult struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID   `bson:"user_id" json:"user_id"`
	Jobs        []primitive.ObjectID `bson:"jobs" json:"jobs"`
	QueryTitle  string               `bson:"query_title" json:"query_title"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
}

type StarredJob struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	JobID  primitive.ObjectID `bson:"job_id" json:"job_id"`
}

type AppliedJob struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	JobID  primitive.ObjectID `bson:"job_id" json:"job_id"`
}

// legacy format for job results for backward compatibility
type JobPageResult struct {
	BusinessName string `json:"business_name"`
	URL          string `json:"url"`
	Description  string `json:"description"`
}
