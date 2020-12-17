package thread

import (
	"cs268-project2-api/graph/model"
	"cs268-project2-api/mongo"
	token2 "cs268-project2-api/token"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func GetThread(input model.ThreadInput) (model.ThreadSingle, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	thread := model.ThreadSingle{
		ID:          "",
		Name:        "",
		TagLine:     "",
		ClassPrefix: "",
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
			return thread, errors
		}
		isValid, _ := token2.ValidateToken(input.Token, userToken, claims["email"].(string))
		fmt.Println(isValid)
		if isValid {
			returnThread, errors := mongo.GetThread(input.ID)
			return returnThread, errors
		} else {
			errors.Errors = true
			errors.Code = 401
			errors.Message = "User not logged in!"
			return thread, errors
		}
	} else {
		errors.Errors = true
		errors.Code = 401
		errors.Message = "User not logged in!"
		return thread, errors
	}
}

func GetThreads(input model.ThreadsInput) ([]*model.Thread, model.Error) {
	errors := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	threads := []*model.Thread{}
	if token, _ := jwt.Parse(input.Token, nil); token != nil {
		// We can parse and extract claims
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims["userID"].(string))
		userToken, returnedErrorGetUser := mongo.GetUserToken(claims["userID"].(string))
		if returnedErrorGetUser.Errors {
			errors.Errors = true
			errors.Code = 500
			errors.Message = "Error in grabbing token"

			return threads, errors
		}
		isValid, _ := token2.ValidateToken(input.Token, userToken, claims["email"].(string))
		fmt.Println(isValid)
		if isValid {
			returnThreads, errors := mongo.GetThreads()
			return returnThreads, errors
		} else {
			errors.Errors = true
			errors.Code = 401
			errors.Message = "User not logged in!"
			return threads, errors
		}
	} else {
		errors.Errors = true
		errors.Code = 401
		errors.Message = "User not logged in!"
		return threads, errors
	}
}
