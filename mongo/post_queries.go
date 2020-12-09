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

type NewPost struct {
	ID        string           `bson:"_id" json:"_id"`
	UserID    string           `json:"userId"`
	Tags      []string         `json:"tags"`
	Content   string           `json:"content"`
	Comments  []*model.Comment `json:"comments"`
	ThreadID  string           `json:"threadId"`
	SubHeader string           `json:"subHeader"`
	Title     string           `json:"title"`
	Class     string           `json:"class"`
	Error     *model.Error     `json:"error"`
}

func GetPost(postID string) model.Post {
	var tags []string
	var comments []*model.Comment
	var errors model.Error
	post := model.Post{
		ID:       "",
		UserID:   "",
		Tags:     tags,
		Content:  "",
		Comments: comments,
		ThreadID: "",
		Error:    &errors,
	}
	err := godotenv.Load(".env")
	if err != nil {
		post.Error.Errors = true
		post.Error.Message = "Failed to load .env"
		post.Error.Code = 500
		return post
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
		post.Error.Errors = true
		post.Error.Message = "Failed to connect to mongo"
		post.Error.Code = 500
		return post
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			post.Error.Errors = true
			post.Error.Message = "Failed to disconnect from mongo"
			post.Error.Code = 500
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		post.Error.Errors = true
		post.Error.Message = "Failed to ping mongo"
		post.Error.Code = 500
		return post
	}

	collection := client.Database("uwecforum").Collection("posts")
	filter := bson.M{
		"$and": []interface{}{
			bson.M{"_id": postID}}}

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	err = collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		post.Error.Errors = true
		post.Error.Message = "Document not found"
		post.Error.Code = 500
		return post
	}
	return post
}

func GetPostsByUser(userID string) model.Post {
	var tags []string
	var comments []*model.Comment
	var errors model.Error
	post := model.Post{
		ID:        "",
		UserID:    "",
		Tags:      tags,
		Content:   "",
		Comments:  comments,
		ThreadID:  "",
		SubHeader: "",
		Title:     "",
		Class:     "",
		Error:     &errors,
	}
	err := godotenv.Load(".env")
	if err != nil {
		post.Error.Errors = true
		post.Error.Message = "Failed to load .env"
		post.Error.Code = 500
		return post
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
		post.Error.Errors = true
		post.Error.Message = "Failed to connect to mongo"
		post.Error.Code = 500
		return post
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			post.Error.Errors = true
			post.Error.Message = "Failed to disconnect from mongo"
			post.Error.Code = 500
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		post.Error.Errors = true
		post.Error.Message = "Failed to ping mongo"
		post.Error.Code = 500
		return post
	}

	collection := client.Database("uwecforum").Collection("posts")
	filter := bson.M{
		"$and": []interface{}{
			bson.M{"userID": userID}}}

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	err = collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		post.Error.Errors = true
		post.Error.Message = "Document not found"
		post.Error.Code = 500
		return post
	}
	return post
}

func GetPostsByThread(threadID string) ([]model.Post, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	var posts []model.Post
	err := godotenv.Load(".env")
	if err != nil {
		errors.Errors = true
		errors.Message = "Failed to load .env"
		errors.Code = 500
		return posts, errors
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
		return posts, errors
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
		return posts, errors
	}

	collection := client.Database("uwecforum").Collection("posts")
	filter := bson.M{
		"$and": []interface{}{
			bson.M{"threadID": threadID}}}

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	cur, err := collection.Find(ctx, filter)
	if err = cur.All(ctx, &posts); err != nil {
		errors.Errors = true
		errors.Message = "Documents not found"
		errors.Code = 500
		return posts, errors
	}
	return posts, errors
}

func CreatePost(input model.CreatePost, userID string) model.Post {
	fmt.Println("Here we are")
	var tags []string
	var comments []*model.Comment
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	idHash, err := bcrypt.GenerateFromPassword([]byte(input.Title), 10)
	post := model.Post{
		UserID:    userID,
		Tags:      tags,
		Content:   input.Content,
		Comments:  comments,
		ThreadID:  input.ThreadID,
		SubHeader: input.SubHeader,
		Title:     input.Title,
		Class:     input.Class,
		Error:     &errors,
		ID:        string(idHash),
	}
	postUpsert := NewPost{
		UserID:    userID,
		Tags:      tags,
		Content:   input.Content,
		Comments:  comments,
		ThreadID:  input.ThreadID,
		SubHeader: input.SubHeader,
		Title:     input.Title,
		Class:     input.Class,
		Error:     &errors,
		ID:        string(idHash),
	}
	err = godotenv.Load(".env")
	if err != nil {
		post.Error.Errors = true
		post.Error.Message = "Failed to load .env"
		post.Error.Code = 500
		return post
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
		post.Error.Errors = true
		post.Error.Message = "Failed to connect to mongo"
		post.Error.Code = 500
		return post
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			post.Error.Errors = true
			post.Error.Message = "Failed to disconnect from mongo"
			post.Error.Code = 500
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		post.Error.Errors = true
		post.Error.Message = "Failed to ping mongo"
		post.Error.Code = 500
		return post
	}

	collection := client.Database("uwecforum").Collection("posts")

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	fmt.Println("Here we are inserting")
	_, err = collection.InsertOne(ctx, postUpsert)
	if err != nil {
		post.Error.Errors = true
		post.Error.Message = "Document not found"
		post.Error.Code = 500
		return post
	}
	fmt.Println(post.ID)
	_, errorsUpdate := UpdateUserPosts(userID, post.ID)
	fmt.Println("Posts should have been updated")
	if errorsUpdate.Errors {
		post.Error = &errorsUpdate
		return post
	}
	return post
}

func EditPost() {
	//TODO: Implement
}

func DeletePost() {
	//TODO: Implement
}
