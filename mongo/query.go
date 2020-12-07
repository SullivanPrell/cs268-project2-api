package mongo

import (
	"context"
	"cs268-project2-api/auth"
	"cs268-project2-api/graph/model"
	validUser "cs268-project2-api/user"
	"fmt"
	"os"
	"os/user"
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
}

func FindOneUser(id string) model.User {
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

	filter := bson.M{
		"$and": []interface{}{
			bson.M{"_id": id}}}

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

func CreateUser(validatedUser validUser.ValidUser) (model.User, model.Error) {
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
	token, errors := auth.GenToken(validatedUser.Email)
	if errors.Errors == true {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Error when generating token"
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
		Token: token,
	}
	_, err = collection.InsertOne(ctx, newUser)
	if err != nil {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Error when creating user"
		errors.Code = 500
		return returnLogin, errors
	}
	user := FindOneUser(validatedUser.ID)

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
	err = collection.FindOne(ctx, filter).Decode(&errors)
	if err != nil {
		fmt.Println(err)
		errors.Errors = true
		errors.Message = "Error when finding user"
		errors.Code = 500
		return false, errors
	}

}
