package mongo

import (
	"context"
	"cs268-project2-api/graph/model"
	token2 "cs268-project2-api/token"
	"cs268-project2-api/validator"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type NewUser struct {
	ID            string              `bson:"_id" json:"_id"`
	Email         string              `json:"email"`
	FirstName     string              `json:"firstName"`
	LastName      string              `json:"lastName"`
	DateOfBirth   string              `json:"dateOfBirth"`
	Major         string              `json:"major"`
	Minor         string              `json:"minor"`
	WillingToHelp bool                `json:"willingToHelp"`
	PostIds       []string            `json:"postIds"`
	CommentIds    []string            `json:"commentIds"`
	ClassesTaken  []string            `json:"classesTaken"`
	EmailVerified model.EmailVerified `json:"emailVerified"`
	Token         model.UserToken     `json:"token"`
	Password      string              `json:"password"`
}

func FindOneUser(idOrEmail string, usingID bool) model.User {
	user := model.User{}
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		user.Error.Errors = true
		user.Error.Message = "Unable to load environment variable"
		user.Error.Code = 500
		return user
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
		user.Error.Errors = true
		user.Error.Message = "Unable to connect to mongo instance"
		user.Error.Code = 503
		return user
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Println(err)
			user.Error.Errors = true
			user.Error.Message = "Unable to disconnect from mongo instance"
			user.Error.Code = 503
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println(err)
		user.Error.Errors = true
		user.Error.Message = "Unable to ping mongo instance"
		user.Error.Code = 503
		return user
	}

	collection := client.Database("uwecforum").Collection("users")
	filter := bson.M{}
	if usingID {
		filter = bson.M{
			"$and": []interface{}{
				bson.M{"_id": idOrEmail}}}
	} else {
		filter = bson.M{
			"$and": []interface{}{
				bson.M{"email": idOrEmail}}}
	}

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		fmt.Println(err)
		user.Error.Errors = true
		user.Error.Message = "Error when finding user"
		user.Error.Code = 500
		return user
	}
	return user
}

func CreateUser(validatedUser validator.ValidUser) (model.User, model.Error) {
	errors := model.Error{}
	returnLogin := model.User{}
	err := godotenv.Load(".env")
	if err != nil {
		errors.Errors = true
		errors.Message = "Unable to load environment variable"
		errors.Code = 500
		return returnLogin, errors
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
		return returnLogin, errors
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
		return returnLogin, errors
	}

	collection := client.Database("uwecforum").Collection("users")

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	token, errors := token2.GenToken(validatedUser.Email, validatedUser.ID)
	if errors.Errors == true {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Error when generating token"
		errors.Code = 500
		return returnLogin, errors
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(validatedUser.Password), 6)
	if err != nil {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Error when hashing password"
		errors.Code = 500
		return returnLogin, errors
	}
	newUser := NewUser{
		FirstName:     validatedUser.FName,
		LastName:      validatedUser.LName,
		ID:            validatedUser.ID,
		Email:         validatedUser.Email,
		DateOfBirth:   validatedUser.DateOfBirth,
		Major:         validatedUser.Major,
		Minor:         validatedUser.Minor,
		WillingToHelp: validatedUser.WillingToHelp,
		PostIds:       []string{},
		CommentIds:    []string{},
		ClassesTaken:  []string{},
		EmailVerified: model.EmailVerified{
			Verified:      false,
			DateValidated: "",
			Email:         validatedUser.Email,
			Error: &model.Error{
				Errors:  false,
				Code:    0,
				Message: "",
			},
		},
		Token:    token,
		Password: string(hash),
	}
	_, err = collection.InsertOne(ctx, newUser)
	if err != nil {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Error when creating user"
		errors.Code = 500
		return returnLogin, errors
	}
	user := FindOneUser(validatedUser.ID, true)

	return user, errors
}

func UserExists(email string) (bool, model.Error) {
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

	collection := client.Database("uwecforum").Collection("users")

	filter := bson.M{
		"$and": []interface{}{
			bson.M{"email": email}}}

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Error when finding user"
		errors.Code = 500
		return false, errors
	}
	if count == 1 {
		return true, errors
	} else if count > 1 {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Error when finding user - multiple users found"
		errors.Code = 500
		return false, errors
	} else {
		fmt.Println(err)
		errors.Errors = false
		errors.Message = "User does not exist"
		errors.Code = 403
		return false, errors
	}
}

func FindUserPass(idOrEmail string, usingID bool) (string, string, model.Error) {
	errors := model.Error{}
	user := NewUser{}
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Unable to load environment variable"
		errors.Code = 500
		return "", "", errors
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
		return "", "", errors
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
		return "", "", errors
	}

	collection := client.Database("uwecforum").Collection("users")
	filter := bson.M{}
	if usingID {
		filter = bson.M{
			"$and": []interface{}{
				bson.M{"_id": idOrEmail}}}
	} else {
		filter = bson.M{
			"$and": []interface{}{
				bson.M{"email": idOrEmail}}}
	}

	ctx, cancelFilter := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFilter()
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Error when finding user"
		errors.Code = 500
		return "", "", errors
	}
	return user.ID, user.Password, errors
}

func UpdateUserToken(id string, token model.UserToken) model.Error {
	errors := model.Error{}
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Unable to load environment variable"
		errors.Code = 500
		return errors
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
		return errors
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
		return errors
	}

	collection := client.Database("uwecforum").Collection("users")
	filter := bson.M{
		"$and": []interface{}{
			bson.M{"_id": id}}}

	update := bson.M{
		"$set": bson.M{"token": token},
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
		errors.Message = "Token update failed"
		errors.Code = 500
		return errors
	}
	return errors
}

func GetUserToken(id string) (model.UserToken, model.Error) {
	errors := model.Error{}
	user := model.User{}
	token := model.UserToken{}
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Unable to load environment variable"
		errors.Code = 500
		return token, errors
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
		return token, errors
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
		return token, errors
	}

	collection := client.Database("uwecforum").Collection("users")
	filter := bson.M{
		"$and": []interface{}{
			bson.M{"_id": id}}}

	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		errors.Errors = true
		errors.Message = "Token update failed"
		errors.Code = 500
		return token, errors
	}

	return *user.Token, errors
}

func UpdateUserPosts(userID string, postID string) (bool, model.Error) {
	fmt.Println("Entered update")
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
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

	collection := client.Database("uwecforum").Collection("users")
	filter := bson.M{
		"$and": []interface{}{
			bson.M{"_id": userID}}}

	update := bson.M{
		"$push": bson.M{"postids": postID},
	}

	// 8) Find one result and update it
	_, err = collection.UpdateOne(ctx, filter, update)
	fmt.Println("Should have executed")
	if err != nil {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Post update failed"
		errors.Code = 500
		return false, errors
	}
	return true, errors
}

func UpdateUserComments(userID string, commentID string) (bool, model.Error) {
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

	collection := client.Database("uwecforum").Collection("users")
	filter := bson.M{
		"$and": []interface{}{
			bson.M{"_id": userID}}}

	update := bson.M{
		"push": bson.M{"commentids": commentID},
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