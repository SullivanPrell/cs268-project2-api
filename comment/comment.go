package comment

import (
	"cs268-project2-api/graph/model"
	"cs268-project2-api/mongo"
	token2 "cs268-project2-api/token"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func CreateNewComment(input model.CreateComment) (model.PostSingle, model.Error) {
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
			postID, errors := mongo.CreateComment(input, claims["userID"].(string))
			if errors.Errors {
				return post, errors
			}
			returnPost, errors := mongo.GetPost(postID)
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

func GetComments(input model.CommentsInput) ([]*model.Comment, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	comments := []*model.Comment{}
	if token, _ := jwt.Parse(input.Token, nil); token != nil {
		// We can parse and extract claims
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims["userID"].(string))
		userToken, returnedErrorGetUser := mongo.GetUserToken(claims["userID"].(string))
		if returnedErrorGetUser.Errors {
			errors.Errors = true
			errors.Code = 500
			errors.Message = "Error in grabbing token"
			return comments, errors
		}
		isValid, _ := token2.ValidateToken(input.Token, userToken, claims["email"].(string))
		fmt.Println(isValid)
		if isValid {
			returnComments, errors := mongo.GetComments(input.CommentIds)
			return returnComments, errors
		} else {
			errors.Errors = true
			errors.Code = 401
			errors.Message = "User not logged in!"
			return comments, errors
		}
	} else {
		errors.Errors = true
		errors.Code = 401
		errors.Message = "User not logged in!"
		return comments, errors
	}
}

func GetCommentsByPostID(input model.CommentsByPostIDInput) ([]*model.Comment, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	comments := []*model.Comment{}
	if token, _ := jwt.Parse(input.Token, nil); token != nil {
		// We can parse and extract claims
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims["userID"].(string))
		userToken, returnedErrorGetUser := mongo.GetUserToken(claims["userID"].(string))
		if returnedErrorGetUser.Errors {
			errors.Errors = true
			errors.Code = 500
			errors.Message = "Error in grabbing token"
			return comments, errors
		}
		isValid, _ := token2.ValidateToken(input.Token, userToken, claims["email"].(string))
		fmt.Println(isValid)
		if isValid {
			returnComments, errors := mongo.GetCommentsByPostID(input.PostID)
			return returnComments, errors
		} else {
			errors.Errors = true
			errors.Code = 401
			errors.Message = "User not logged in!"
			return comments, errors
		}
	} else {
		errors.Errors = true
		errors.Code = 401
		errors.Message = "User not logged in!"
		return comments, errors
	}
}
