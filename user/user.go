package user

import (
	"cs268-project2-api/graph/model"
	"cs268-project2-api/mongo"
	token2 "cs268-project2-api/token"
	"cs268-project2-api/validator"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func CreateNewUser(input model.CreateUser) model.User {
	blankError := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	blankToken := model.UserToken{
		Token:      "",
		ExpireDate: 0,
	}
	blankEmailVerified := model.EmailVerified{
		Verified:      false,
		DateValidated: "",
		Email:         "",
		Error:         &blankError,
	}
	blankUser := model.User{
		ID:            "",
		Email:         "",
		FirstName:     "",
		LastName:      "",
		DateOfBirth:   "",
		Major:         "",
		Minor:         "",
		WillingToHelp: false,
		PostIds:       []string{},
		CommentIds:    []string{},
		ClassesTaken:  []string{},
		EmailVerified: &blankEmailVerified,
		Token:         &blankToken,
		Error:         &blankError,
	}
	userExists, errors := mongo.UserExists(input.Email)
	if errors.Errors == true {
		fmt.Print(errors.Message)
		blankUser.Error = &errors
		return blankUser
	}
	if userExists {
		userExistsError := model.Error{
			Errors:  true,
			Code:    403,
			Message: "User already exists!",
		}
		blankUser.Error = &userExistsError
		return blankUser
	}
	validatedUser, errors := validator.ValidateInfo(input)
	if errors.Errors == true {
		fmt.Print(errors.Message)
		blankUser.Error = &errors
		return blankUser
	}
	user, errors := mongo.CreateUser(validatedUser)
	if errors.Errors == true {
		fmt.Print(errors.Message)
		blankUser.Error = &errors
		return blankUser
	}
	user.Error = &errors
	return user
}

func FindUser(input model.UserInput) model.User {
	returnError := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	blankError := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	blankToken := model.UserToken{
		Token:      "",
		ExpireDate: 0,
	}
	blankEmailVerified := model.EmailVerified{
		Verified:      false,
		DateValidated: "",
		Email:         "",
		Error:         &blankError,
	}
	loginReturn := model.User{
		ID:            "",
		Email:         "",
		FirstName:     "",
		LastName:      "",
		DateOfBirth:   "",
		Major:         "",
		Minor:         "",
		WillingToHelp: false,
		PostIds:       []string{},
		CommentIds:    []string{},
		ClassesTaken:  []string{},
		EmailVerified: &blankEmailVerified,
		Token:         &blankToken,
		Error:         &blankError,
	}
	if token, _ := jwt.Parse(input.Token, nil); token != nil {

		// We can parse and extract claims
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims["userID"].(string))
		userToken, returnedErrorGetUser := mongo.GetUserToken(claims["userID"].(string))
		if returnedErrorGetUser.Errors {
			loginReturn.Error = &returnedErrorGetUser
			return loginReturn
		}
		isValid, _ := token2.ValidateToken(input.Token, userToken, claims["email"].(string))
		fmt.Println(isValid)
		if isValid {
			returnUser := mongo.FindOneUser(claims["userID"].(string), true)
			return returnUser
		} else {
			returnError.Errors = true
			returnError.Code = 401
			returnError.Message = "User not logged in!"
			loginReturn.Error = &returnError
			return loginReturn
		}
	} else {
		returnError.Errors = true
		returnError.Code = 500
		returnError.Message = "Token invalid"
		loginReturn.Error = &returnError
		return loginReturn
	}
}
