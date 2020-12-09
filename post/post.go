package post

import (
	"cs268-project2-api/graph/model"
	"cs268-project2-api/mongo"
	token2 "cs268-project2-api/token"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func CreateNewPost(input model.CreatePost) model.Post {
	var tags []string
	var comments []*model.Comment
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	post := model.Post{
		ID:       "",
		UserID:   "",
		Tags:     tags,
		Content:  "",
		Comments: comments,
		ThreadID: "",
		Error:    &errors,
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
			post.Error = &returnedErrorGetUser
			return post
		}
		isValid, _ := token2.ValidateToken(input.Token, userToken, claims["email"].(string))
		fmt.Println(isValid)
		if isValid {
			returnPost := mongo.CreatePost(input, claims["userID"].(string))
			return returnPost
		} else {
			errors.Errors = true
			errors.Code = 401
			errors.Message = "User not logged in!"
			post.Error = &errors
			return post
		}
	} else {
		errors.Errors = true
		errors.Code = 401
		errors.Message = "User not logged in!"
		post.Error = &errors
		return post
	}
}