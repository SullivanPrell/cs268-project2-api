package post

import (
	"cs268-project2-api/graph/model"
	"cs268-project2-api/mongo"
	token2 "cs268-project2-api/token"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func CreateNewPost(input model.CreatePost) (model.PostSingle, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	post := model.PostSingle{
		ID:         "",
		UserID:     "",
		Content:    "",
		CommentIDs: []string{},
		ThreadID:   "",
	}
	if token, _ := jwt.Parse(input.Token, nil); token != nil {
		// We can parse and extract claims
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims["userID"].(string))
		userToken, returnedErrorGetUser := mongo.GetUserToken(claims["userID"].(string))
		if returnedErrorGetUser.Errors {
			errors.Errors = true
			errors.Code = 500
			errors.Message = "Error in grabbing token"

			return post, errors
		}
		isValid, _ := token2.ValidateToken(input.Token, userToken, claims["email"].(string))
		fmt.Println(isValid)
		if isValid {
			returnPost, errors := mongo.CreatePost(input, claims["userID"].(string))
			return returnPost, errors
		} else {
			errors.Errors = true
			errors.Code = 401
			errors.Message = "User not logged in!"

			return post, errors
		}
	} else {
		errors.Errors = true
		errors.Code = 401
		errors.Message = "User not logged in!"
		return post, errors
	}
}

func GetPost(input model.PostInput) (model.PostSingle, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	post := model.PostSingle{
		ID:         "",
		UserID:     "",
		Content:    "",
		CommentIDs: []string{},
		ThreadID:   "",
	}
	if token, _ := jwt.Parse(input.Token, nil); token != nil {
		// We can parse and extract claims
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims["userID"].(string))
		userToken, returnedErrorGetUser := mongo.GetUserToken(claims["userID"].(string))
		if returnedErrorGetUser.Errors {
			errors.Errors = true
			errors.Code = 500
			errors.Message = "Error in grabbing token"

			return post, errors
		}
		isValid, _ := token2.ValidateToken(input.Token, userToken, claims["email"].(string))
		fmt.Println(isValid)
		if isValid {
			returnPost, errors := mongo.GetPost(input.PostID)
			return returnPost, errors
		} else {
			errors.Errors = true
			errors.Code = 401
			errors.Message = "User not logged in!"
			return post, errors
		}
	} else {
		errors.Errors = true
		errors.Code = 401
		errors.Message = "User not logged in!"
		return post, errors
	}
}

func GetPosts(input model.PostsInput) ([]*model.Post, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	posts := []*model.Post{}
	if token, _ := jwt.Parse(input.Token, nil); token != nil {
		// We can parse and extract claims
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims["userID"].(string))
		userToken, returnedErrorGetUser := mongo.GetUserToken(claims["userID"].(string))
		if returnedErrorGetUser.Errors {
			errors.Errors = true
			errors.Code = 500
			errors.Message = "Error in grabbing token"
			return posts, errors
		}
		isValid, _ := token2.ValidateToken(input.Token, userToken, claims["email"].(string))
		fmt.Println(isValid)
		if isValid {
			returnPosts, errors := mongo.GetPostsByPostIDs(input.PostIds)
			return returnPosts, errors
		} else {
			errors.Errors = true
			errors.Code = 401
			errors.Message = "User not logged in!"
			return posts, errors
		}
	} else {
		errors.Errors = true
		errors.Code = 401
		errors.Message = "User not logged in!"
		return posts, errors
	}
}

func GetPostsByUserID(input model.PostsByUserIDInput) ([]*model.Post, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	posts := []*model.Post{}
	if token, _ := jwt.Parse(input.Token, nil); token != nil {
		// We can parse and extract claims
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims["userID"].(string))
		userToken, returnedErrorGetUser := mongo.GetUserToken(claims["userID"].(string))
		if returnedErrorGetUser.Errors {
			errors.Errors = true
			errors.Code = 500
			errors.Message = "Error in grabbing token"
			return posts, errors
		}
		isValid, _ := token2.ValidateToken(input.Token, userToken, claims["email"].(string))
		fmt.Println(isValid)
		if isValid {
			returnPosts, errors := mongo.GetPostsByUser(input.UserID)
			return returnPosts, errors
		} else {
			errors.Errors = true
			errors.Code = 401
			errors.Message = "User not logged in!"
			return posts, errors
		}
	} else {
		errors.Errors = true
		errors.Code = 401
		errors.Message = "User not logged in!"
		return posts, errors
	}
}

func GetPostsByThread(input model.PostsByThreadInput) ([]*model.Post, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	posts := []*model.Post{}
	if token, _ := jwt.Parse(input.Token, nil); token != nil {
		// We can parse and extract claims
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims["userID"].(string))
		userToken, returnedErrorGetUser := mongo.GetUserToken(claims["userID"].(string))
		if returnedErrorGetUser.Errors {
			errors.Errors = true
			errors.Code = 500
			errors.Message = "Error in grabbing token"
			return posts, errors
		}
		isValid, _ := token2.ValidateToken(input.Token, userToken, claims["email"].(string))
		fmt.Println(isValid)
		if isValid {
			returnPosts, errors := mongo.GetPostsByThread(input.ThreadID)
			return returnPosts, errors
		} else {
			errors.Errors = true
			errors.Code = 401
			errors.Message = "User not logged in!"
			return posts, errors
		}
	} else {
		errors.Errors = true
		errors.Code = 401
		errors.Message = "User not logged in!"
		return posts, errors
	}
}
