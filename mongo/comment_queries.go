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
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type NewComment struct {
	ID      string `bson:"_id" json:"_id"`
	UserID  string `json:"userId"`
	Content string `json:"content"`
	PostID  string `json:"postId"`
}

func CreateComment(input model.CreateComment, userID string) (string, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	idHash, err := bcrypt.GenerateFromPassword([]byte(input.PostID), 4)

	commentUpsert := NewComment{
		UserID:  userID,
		Content: input.Content,
		ID:      string(idHash),
		PostID:  input.PostID,
	}
	err = godotenv.Load(".env")
	if err != nil {
		errors.Errors = true
		errors.Message = "Failed to load .env"
		errors.Code = 500
		return input.PostID, errors
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
		return input.PostID, errors
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
		return input.PostID, errors
	}

	collection := client.Database("uwecforum").Collection("comments")

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	fmt.Println("Here we are inserting")
	_, err = collection.InsertOne(ctx, commentUpsert)
	if err != nil {
		errors.Errors = true
		errors.Message = "Document not found"
		errors.Code = 500
		return input.PostID, errors
	}
	_, errorsUpdate := UpdateUserComments(userID, commentUpsert.ID)
	fmt.Println("Posts should have been updated")
	if errorsUpdate.Errors {
		errors = errorsUpdate
		return input.PostID, errors
	}
	_, errorsUpdate = UpdatePostComments(input.PostID, string(idHash))
	return input.PostID, errors
}

func GetCommentsByPostID(postID string) ([]*model.Comment, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	var comments []*model.Comment
	err := godotenv.Load(".env")
	if err != nil {
		errors.Errors = true
		errors.Message = "Failed to load .env"
		errors.Code = 500
		return comments, errors
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
		return comments, errors
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
		return comments, errors
	}

	collection := client.Database("uwecforum").Collection("comments")
	filter := bson.M{
		"$and": []interface{}{
			bson.M{"postid": postID}}}

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	cur, err := collection.Find(ctx, filter)
	if err = cur.All(ctx, &comments); err != nil {
		errors.Errors = true
		errors.Message = "Documents not found"
		errors.Code = 500
		return comments, errors
	}
	return comments, errors
}

func GetComments(commentIDs []string) ([]*model.Comment, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	var comments []*model.Comment
	err := godotenv.Load(".env")
	if err != nil {
		errors.Errors = true
		errors.Message = "Failed to load .env"
		errors.Code = 500
		return comments, errors
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
		return comments, errors
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
		return comments, errors
	}

	collection := client.Database("uwecforum").Collection("comments")
	filter := bson.M{
		"$in": []interface{}{
			bson.M{"_id": commentIDs}}}

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	cur, err := collection.Find(ctx, filter)
	if err = cur.All(ctx, &comments); err != nil {
		errors.Errors = true
		errors.Message = "Documents not found"
		errors.Code = 500
		return comments, errors
	}
	return comments, errors
}
