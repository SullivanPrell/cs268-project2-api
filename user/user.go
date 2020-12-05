package user

import (
	"cs268-project2-api/graph/model"

	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/validator.v2"
)

type UserValidator struct {
	FName         string `validate:"nonzero,min=2,max=100"`
	LName         string `validate:"nonzero,min=2,max=100"`
	Email         string `validate:"nonzero"`
	Password      string `validate:"min=10,max=350"`
	ID            string `validate:"nonzero"`
	Major         string `validate:"nonzero"`
	Minor         string `validate:"nonzero"`
	DateOfBirth   string `validate:"nonzero"`
	WillingToHelp bool   `validate:"nonzero"`
}

type Errors struct {
	errors    bool
	errorList []string
}

type ValidUser struct {
	Email         string
	Password      string
	PasswordHash  string
	DateOfBirth   string
	Major         string
	Minor         string
	WillingToHelp bool
	ID            string
	FName         string
	LName         string
}

// TODO: Fix errors

func ValidateInfo(apiUser model.CreateUser) (ValidUser, model.Error) {
	errors := model.Error{}
	returnUser := ValidUser{}
	idHash, err := bcrypt.GenerateFromPassword([]byte(apiUser.Email), 10)
	if err != nil {
		errors.Errors = true
		errors.Message = "Failed to encrypt ID"
		errors.Code = 500
		return returnUser, errors
	}
	validateUser := UserValidator{
		Email:         apiUser.Email,
		FName:         apiUser.FirstName,
		LName:         apiUser.LastName,
		Password:      apiUser.Password,
		Major:         apiUser.Major,
		Minor:         *apiUser.Minor,
		DateOfBirth:   apiUser.DateOfBirth,
		WillingToHelp: apiUser.WillingToHelp,
		ID:            string(idHash),
	}

	validateUser.ID = string(idHash)

	err = validator.Validate(validateUser)
	if err != nil {
		errors.Errors = true
		errors.Message = "Error in validation of user info. Please follow website guidelines."
		errors.Code = 400
		return returnUser, errors
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(apiUser.Password), 12)
	if err != nil {
		errors.Errors = true
		errors.Message = "Error when hashing password."
		errors.Code = 500
		return returnUser, errors
	}

	returnUser = ValidUser{
		Email:         apiUser.Email,
		FName:         apiUser.FirstName,
		LName:         apiUser.LastName,
		Password:      apiUser.Password,
		PasswordHash:  string(passHash),
		Major:         apiUser.Major,
		Minor:         *apiUser.Minor,
		DateOfBirth:   apiUser.DateOfBirth,
		WillingToHelp: apiUser.WillingToHelp,
		ID:            string(idHash),
	}
	errors.Errors = false

	return returnUser, errors
}

// func NewUser(userInfo ValidatedUser, collection *gocb.Collection) (UserToken, MutationPayload) {
// 	returnErr := MutationPayload{}
// 	returnErr.Success = true
// 	returnToken := userInfo.ValidUser.Token
// 	if userInfo.UserValid {
// 		exists, _ := UserExist(userInfo.ValidUser.Email, collection)
// 		if exists {
// 			returnErr.Success = false
// 			returnErr.Errors = append(returnErr.Errors, "Email already in use!")
// 			returnErr.Token = ""
// 			returnToken.Token = ""
// 			returnToken.ExpireDate = 0000

// 		} else {
// 			_, err := collection.Upsert(userInfo.ValidUser.Email, userInfo.ValidUser, &gocb.UpsertOptions{})
// 			if err != nil {
// 				returnErr.Success = false
// 				returnErr.Errors = append(returnErr.Errors, "Account Creation Error, please try again later. Dev Code: ERRNEWUSRDBUP")
// 				returnErr.Token = ""
// 				returnToken.Token = ""
// 				returnToken.ExpireDate = 0000
// 			}

// 		}

// 	} else {
// 		returnErr.Errors = append(returnErr.Errors, userInfo.Errors...)
// 		returnErr.Success = false
// 		returnErr.Errors = append(returnErr.Errors, "Account Creation Error, please try again later. Dev Code: ERRNEWUSRNVU")
// 		returnErr.Token = ""
// 		returnToken.Token = ""
// 		returnToken.ExpireDate = 0000
// 		return returnToken, returnErr
// 	}

// 	return returnToken, returnErr
// }

// func UpdateUser(modifyDetails ModifyUser, collection *gocb.Collection) bool {
// 	//TODO: Implement
// 	return false
// }

// func RemoveUser(userInfo User) bool {
// 	//TODO: Implement
// 	return false
// }

// func UserExist(email string, collection *gocb.Collection) (bool, APIError) {
// 	apiErr := APIError{}
// 	checkUser, _ := collection.Get(email, nil)
// 	//TODO: Swap to check exist function couchbase
// 	// if err != nil {
// 	// 	if err.error_name == "KEY_ENOENT" {
// 	// 		return false, apiErr
// 	// 	}
// 	// 	apiErr.Error = true
// 	// 	apiErr.Message = "Account Validation Error, please try again later."
// 	// 	panic(err)
// 	// 	return false, apiErr
// 	// }
// 	if checkUser != nil {
// 		return true, apiErr
// 	}
// 	return false, apiErr

// }
