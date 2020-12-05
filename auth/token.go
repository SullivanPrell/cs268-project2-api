package auth

import (
	"cs268-project2-api/graph/model"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	jwt "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenToken(email string) (model.UserToken, model.Error) {
	newError := model.Error{}
	err := godotenv.Load(".env")
	if err != nil {
		newError.Errors = true
		newError.Code = 500
		newError.Message = "Error reading from env file"
	}
	jwtSecret := fmt.Sprintf("%s%s", os.Getenv("JWT_SECRET"), email)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(jwtSecret))
	returnUserToken := model.UserToken{
		Token:     tokenString,
		ExpireDate: int(expirationTime.Unix()),
	}
	if err != nil {
		newError.Errors = true
		newError.Message = "Error when creating token"
		newError.Code = 500
	}
	return returnUserToken, newError
}

func ValidateToken(tokenString string, userTokenInfo model.UserToken, email string) bool {

	// Initialize a new instance of `Claims`
	if tokenString == userTokenInfo.Token {

		// Pass 1 Tokens are equal!
		if time.Now().Unix() > int64(userTokenInfo.ExpireDate) {
			// TOKEN EXPIRED !
			return false
		} else {
			claims := &Claims{}
			//getEnv()
			jwtSecret := fmt.Sprintf("%s%s", os.Getenv("JWT_SECRET"), email)
			// Parse the JWT string and store the result in `claims`.
			// Note that we are passing the key in this method as well. This method will return an error
			// if the token is invalid (if it has expired according to the expiry time we set on sign in),
			// or if the signature does not match
			tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {

					return false
				}
			}
			if !tkn.Valid {

				return false
			}
			return true
		}

	} else {
		return false
	}

}

//func RenewToken(email string, id string, tokenString string, tokenExpire int64) (UserToken, APIError) {
//	renewToken := UserToken{}
//	renewError := APIError{}
//	// If expires in 30 seconds
//	if tokenExpire > (time.Now().Add(time.Hour * 24)).Unix() {
//		// Not expired
//	} else {
//		// Generate a new token, pass token back
//		renewToken, renewError = GenToken(email, id)
//
//	}
//
//	return renewToken, renewError
//
//}
//
//func GetCurrentToken(email string, collection *gocb.Collection) (UserToken, APIError) {
//	returnToken := UserToken{}
//	returnError := APIError{}
//
//	//TODO: Convert to Mongo from Couchbase
//
//	ops := []gocb.LookupInSpec{
//		gocb.GetSpec("token.token", &gocb.GetSpecOptions{}),
//		gocb.GetSpec("token.expiredate", &gocb.GetSpecOptions{}),
//	}
//	getResult, err := collection.LookupIn(email, ops, &gocb.LookupInOptions{})
//	if err != nil {
//		panic(err)
//		//TODO: Create API Err
//	}
//
//	var currentToken string
//	var currentExpireDate int64
//	err = getResult.ContentAt(0, &currentToken)
//	if err != nil {
//		panic(err)
//		// Create API Err
//	}
//	err = getResult.ContentAt(1, &currentExpireDate)
//	if err != nil {
//		panic(err)
//		// Create API Err
//	}
//	returnToken.Token = currentToken
//	returnToken.ExpireDate = currentExpireDate
//	return returnToken, returnError
//}
