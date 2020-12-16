package validator

import (
	"cs268-project2-api/graph/model"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
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
		fmt.Print(err)
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
