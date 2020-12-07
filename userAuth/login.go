package userAuth

import (
	"cs268-project2-api/mongo"
	token2 "cs268-project2-api/token"
	"fmt"
	"github.com/dgrijalva/jwt-go"

	"cs268-project2-api/graph/model"
	"golang.org/x/crypto/bcrypt"
)

func Login(userInput model.LoginInput) model.Login {
	blankToken := model.UserToken{
		Token:      "",
		ExpireDate: 0,
	}
	returnError := model.Error{
		Errors:  false,
		Code:    0,
		Message: "",
	}
	loginReturn := model.Login{
		UserID: "",
		Email:  "",
		Token:  &blankToken,
		Error:  &returnError,
	}

	if userInput.Token == "" && userInput.Email != "" && userInput.Password != "" {
		fmt.Println("Pass Login")
		// Password Login
		id, pass, returnError := mongo.FindUserPass(userInput.Email, false)
		if returnError.Errors {
			loginReturn.Error = &returnError
			return loginReturn
		}

		err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(userInput.Password))

		if err != nil {
			returnError.Errors = true
			returnError.Code = 403
			returnError.Message = "Incorrect Password"
			loginReturn.Error = &returnError
		}

		token, returnError := token2.GenToken(userInput.Email, id)
		if returnError.Errors {
			loginReturn.Error = &returnError
			return loginReturn
		}

		returnError = mongo.UpdateUserToken(id, token)
		if returnError.Errors {
			loginReturn.Error = &returnError
			return loginReturn
		}

		loginReturn.Email = userInput.Email
		loginReturn.UserID = id
		loginReturn.Token = &token
		loginReturn.Error = &returnError
		return loginReturn

	} else if userInput.Token != "" {
		fmt.Println("Token Login")
		// Token Login
		if token, _ := jwt.Parse(userInput.Token, nil); token != nil {
			// We can parse and extract claims
			claims := token.Claims.(jwt.MapClaims)
			fmt.Println(claims["userID"].(string))
			userToken, returnedErrorGetUser := mongo.GetUserToken(claims["userID"].(string))
			if returnedErrorGetUser.Errors {
				loginReturn.Error = &returnedErrorGetUser
				return loginReturn
			}
			isValid, _ := token2.ValidateToken(userInput.Token, userToken, claims["email"].(string))
			fmt.Println(isValid)
			if isValid {
				// Token is valid - login successful
				renewed, newToken, returnedErrorRenew := token2.RenewToken(claims["userID"].(string), claims["email"].(string), userToken)
				fmt.Println(renewed)
				if returnedErrorRenew.Errors {
					loginReturn.Error = &returnedErrorRenew
					return loginReturn
				}
				if renewed {
					// Token has been renewed - need to update - need to update cookies as well
					returnedErrorUpdate := mongo.UpdateUserToken(claims["userID"].(string), newToken)
					if returnedErrorUpdate.Errors {
						loginReturn.Error = &returnedErrorUpdate
						return loginReturn
					}
				}
				returnError = model.Error{
					Errors:  false,
					Code:    0,
					Message: "",
				}
				loginReturn.Token.Token = newToken.Token
				loginReturn.Token.ExpireDate = newToken.ExpireDate
				loginReturn.UserID = claims["userID"].(string)
				loginReturn.Email = claims["email"].(string)
				loginReturn.Error = &returnError
				return loginReturn

			}
		}

	} else {
		fmt.Println("Error Login")
		// Error - input is invalid
		returnError.Errors = true
		returnError.Message = "Missing required fields"
		returnError.Code = 400
		loginReturn.Error = &returnError
		return loginReturn
	}
	return loginReturn
}
