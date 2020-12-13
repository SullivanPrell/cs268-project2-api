package mongo

import (
	"context"
	"cs268-project2-api/graph/model"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

func GetThread(threadID string) (model.ThreadSingle, model.Error) {
	var errors model.Error
	thread := model.ThreadSingle{
		ID:          "",
		Name:        "",
		TagLine:     "",
		ClassPrefix: "",
	}
	err := godotenv.Load(".env")
	if err != nil {
		errors.Errors = true
		errors.Message = "Failed to load .env"
		errors.Code = 500
		return thread, errors
	}
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// Replace the uri string with your MongoDB deployment's connection string.
	uri := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=%s", os.Getenv("MONGO_USR"), os.Getenv("MONGO_PASS"), os.Getenv("MONGO_URL"), os.Getenv("MONGO_AUTHDB"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		errors.Errors = true
		errors.Message = "Failed to connect to mongo"
		errors.Code = 500
		return thread, errors
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			errors.Errors = true
			errors.Message = "Failed to disconnect from mongo"
			errors.Code = 500
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		errors.Errors = true
		errors.Message = "Failed to ping mongo"
		errors.Code = 500
		return thread, errors
	}

	collection := client.Database("uwecforum").Collection("threads")
	filter := bson.M{
		"$and": []interface{}{
			bson.M{"_id": threadID}}}

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	err = collection.FindOne(ctx, filter).Decode(&thread)
	if err != nil {
		errors.Errors = true
		errors.Message = "Document not found"
		errors.Code = 500
		return thread, errors
	}
	return thread, errors
}

func GetThreads() ([]*model.Thread, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	var threads []*model.Thread
	err := godotenv.Load(".env")
	if err != nil {
		errors.Errors = true
		errors.Message = "Failed to load .env"
		errors.Code = 500
		return threads, errors
	}
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// Replace the uri string with your MongoDB deployment's connection string.
	uri := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=%s", os.Getenv("MONGO_USR"), os.Getenv("MONGO_PASS"), os.Getenv("MONGO_URL"), os.Getenv("MONGO_AUTHDB"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		errors.Errors = true
		errors.Message = "Failed to connect to mongo"
		errors.Code = 500
		return threads, errors
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			errors.Errors = true
			errors.Message = "Failed to disconnect from mongo"
			errors.Code = 500
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		errors.Errors = true
		errors.Message = "Failed to ping mongo"
		errors.Code = 500
		return threads, errors
	}

	collection := client.Database("uwecforum").Collection("threads")
	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	cur, err := collection.Find(ctx, bson.M{})
	if err = cur.All(ctx, &threads); err != nil {
		errors.Errors = true
		errors.Message = "Documents not found"
		errors.Code = 500
		return threads, errors
	}
	return threads, errors
}
