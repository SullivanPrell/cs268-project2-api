package tokenGen

import (
	"cs268-project2-api/graph/model"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	jwt "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Email  string `json:"email"`
	UserID string `json:"userID"`
	jwt.StandardClaims
}

func GenToken(email string, userID string) (model.UserToken, model.Error) {
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
		Email:  email,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(jwtSecret))
	returnUserToken := model.UserToken{
		Token:      tokenString,
		ExpireDate: int(expirationTime.Unix()),
	}
	if err != nil {
		newError.Errors = true
		newError.Message = "Error when creating token"
		newError.Code = 500
	}
	return returnUserToken, newError
}

func ValidateToken(tokenString string, userTokenInfo model.UserToken, email string) (bool, jwt.MapClaims) {

	// Initialize a new instance of `Claims`
	if tokenString == userTokenInfo.Token {

		// Pass 1 Tokens are equal!
		if time.Now().Unix() > int64(userTokenInfo.ExpireDate) {
			// TOKEN EXPIRED !
			return false, jwt.MapClaims{}
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

					return false, jwt.MapClaims{}
				}
			}
			if !tkn.Valid {

				return false, jwt.MapClaims{}
			}
			if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
				return true, claims
			}
			return true, jwt.MapClaims{}
		}

	} else {
		return false, jwt.MapClaims{}
	}

}

func RenewToken(email string, userID string, tokenInfo model.UserToken) (bool, model.UserToken, model.Error) {
	renewToken := model.UserToken{}
	renewError := model.Error{}
	// If expires in 4 hours
	if int64(tokenInfo.ExpireDate) > (time.Now().Add(time.Hour * 4)).Unix() {
		// Not expired
		return false, tokenInfo, renewError
	} else {
		// Generate a new token, pass token back
		renewToken, renewError = GenToken(email, userID)
		return true, renewToken, renewError
	}

}
