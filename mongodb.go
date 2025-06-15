package mongostore

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Conn *mongo.Client

func Init(ctx context.Context, dsn string) {
	var err error
	Conn, err = mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		log.Fatalf("failed to connect to mongodb %s", err)
	}

	if err = pingMongoDB(ctx); err != nil {
		log.Fatalf("failed to ping mongodb server %s", err)
	}
	slog.Info("connected to mongodb")
}

func pingMongoDB(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return Conn.Ping(ctx, readpref.Primary())
}

func Shutdown(ctx context.Context) {
	if err := Conn.Disconnect(ctx); err != nil {
		slog.Error("Error disconnecting from mongodb", slog.String("Error", err.Error()))
	}

	Conn = nil
}

type CreateIndexParams struct {
	DatabaseName   string
	CollectionName string
	IsUnique       bool
	Fields         bson.D
}

func CreateIndex(params CreateIndexParams) bool {

	// Define index options
	indexOpts := options.Index().SetUnique(params.IsUnique)

	mod := mongo.IndexModel{
		Keys:    params.Fields,
		Options: indexOpts,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	collection := Conn.Database(params.DatabaseName).Collection(params.CollectionName)
	idxName, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		slog.Info("CreateIndex", "err", fmt.Sprintf("Unable to create index on collection '%s' .", params.CollectionName))

		return false
	}

	slog.Info(fmt.Sprintf("Successfully created index on collection '%s', index name is '%s'", params.CollectionName, idxName))

	return true
}
