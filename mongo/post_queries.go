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
	ID         string   `bson:"_id" json:"_id"`
	UserID     string   `json:"userId"`
	Content    string   `json:"content"`
	CommentIDs []string `json:"commentIDs"`
	ThreadID   string   `json:"threadId"`
	SubHeader  string   `json:"subHeader"`
	Title      string   `json:"title"`
	Class      string   `json:"class"`
}

func GetPost(postID string) (model.PostSingle, model.Error) {
	var errors model.Error
	post := model.PostSingle{
		ID:         "",
		UserID:     "",
		Content:    "",
		CommentIDs: []string{},
		ThreadID:   "",
	}
	err := godotenv.Load(".env")
	if err != nil {
		errors.Errors = true
		errors.Message = "Failed to load .env"
		errors.Code = 500
		return post, errors
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
		return post, errors
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
		return post, errors
	}

	collection := client.Database("uwecforum").Collection("posts")
	filter := bson.M{
		"$and": []interface{}{
			bson.M{"_id": postID}}}

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	err = collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		errors.Errors = true
		errors.Message = "Document not found"
		errors.Code = 500
		return post, errors
	}
	return post, errors
}

func GetPostsByUser(userID string) ([]*model.Post, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	var posts []*model.Post
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
			bson.M{"userid": userID}}}

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

func GetPostsByThread(threadID string) ([]*model.Post, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	var posts []*model.Post
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
			bson.M{"threadid": threadID}}}

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

func GetPostsByPostIDs(postIDs []string) ([]*model.Post, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	var posts []*model.Post
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
		"$in": []interface{}{
			bson.M{"_id": postIDs}}}

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

func CreatePost(input model.CreatePost, userID string) (model.PostSingle, model.Error) {
	fmt.Println("Here we are")
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	idHash, err := bcrypt.GenerateFromPassword([]byte(input.Title), 10)
	post := model.PostSingle{
		UserID:     userID,
		Content:    input.Content,
		CommentIDs: []string{},
		ThreadID:   input.ThreadID,
		SubHeader:  input.SubHeader,
		Title:      input.Title,
		Class:      input.Class,
		ID:         string(idHash),
	}
	postUpsert := NewPost{
		UserID:     userID,
		Content:    input.Content,
		CommentIDs: []string{},
		ThreadID:   input.ThreadID,
		SubHeader:  input.SubHeader,
		Title:      input.Title,
		Class:      input.Class,
		ID:         string(idHash),
	}
	err = godotenv.Load(".env")
	if err != nil {
		errors.Errors = true
		errors.Message = "Failed to load .env"
		errors.Code = 500
		return post, errors
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
		return post, errors
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
		return post, errors
	}

	collection := client.Database("uwecforum").Collection("posts")

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	fmt.Println("Here we are inserting")
	_, err = collection.InsertOne(ctx, postUpsert)
	if err != nil {
		errors.Errors = true
		errors.Message = "Document not found"
		errors.Code = 500
		return post, errors
	}
	fmt.Println(post.ID)
	_, errorsUpdate := UpdateUserPosts(userID, post.ID)
	fmt.Println("Posts should have been updated")
	if errorsUpdate.Errors {
		errors = errorsUpdate
		return post, errors
	}
	return post, errors
}

func UpdatePostComments(postID string, commentID string) (bool, model.Error) {
	errors := model.Error{}
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Unable to load environment variable"
		errors.Code = 500
		return false, errors
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
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Unable to connect to mongo instance"
		errors.Code = 503
		return false, errors
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Println(err)
			errors.Errors = true
			errors.Message = "Unable to disconnect from mongo instance"
			errors.Code = 503
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Unable to ping mongo instance"
		errors.Code = 503
		return false, errors
	}

	collection := client.Database("uwecforum").Collection("posts")
	filter := bson.M{
		"$and": []interface{}{
			bson.M{"_id": postID}}}

	update := bson.M{
		"$push": bson.M{"commentids": commentID},
	}

	// 7) Create an instance of an options and set the desired options
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	// 8) Find one result and update it
	result := collection.FindOneAndUpdate(ctx, filter, update, &opt)
	if result.Err() != nil {
		errors.Errors = true
		errors.Message = "Comment update failed"
		errors.Code = 500
		return false, errors
	}
	return true, errors
}

func EditPost() {
	//TODO: Implement
}

func DeletePost() {
	//TODO: Implement
}
