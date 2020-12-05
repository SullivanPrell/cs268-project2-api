package auth

//import (
//	"time"
//
//	"github.com/couchbase/gocb"
//	"golang.org/x/crypto/bcrypt"
//	"cs268-project2-api/graph/model"
//)

//func Login(model.LoginInput) model.Login {
//	returnError := APIError{}
//	emptyUserReturn := UserToken{}
//	var dbPassword string
//	var dbID string
//	var dbToken string
//	var dbExpireDate int64
//	var dbUserToken UserToken
//	userExists, errors := UserExist(email, collection)
//	if errors.Error {
//		return emptyUserReturn, errors
//	} else if !userExists {
//		returnError.Error = true
//		returnError.Message = "User does not exist!"
//		return emptyUserReturn, returnError
//	} else {
//		// User Exists and no errors!
//
//		dbUserToken = UserToken{
//			Token:      dbToken,
//			ExpireDate: dbExpireDate,
//		}
//		if !ValidateToken(jwt, dbUserToken, dbID) {
//
//			err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
//			if err != nil {
//				returnError.Error = true
//				returnError.Message = "Incorrect password!"
//				return emptyUserReturn, returnError
//			} else {
//				//TODO: Generate and return new token
//				newToken, tokenGenErrors := GenToken(email, dbID)
//				if tokenGenErrors.Error {
//					returnError.Error = true
//					returnError.Message = "Token Generation Error, please try again later. Dev Code: LOGTOKGENERR"
//					return emptyUserReturn, returnError
//				} else {
//					// No error return token and add to DB
//
//					return newToken, returnError
//				}
//			}
//		}
//
//	}
//	return dbUserToken, returnError
//}