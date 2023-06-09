package resolver

import (
	"errors"
	"fmt"
	"lms/database"
	"lms/gosql"
	validation "lms/graphql/validate"
	"lms/model"
	"lms/utils"
	"reflect"

	"github.com/graphql-go/graphql"
)

func Login(params graphql.ResolveParams) (interface{}, error) {
	phone := params.Args["phone"].(string)
	password := params.Args["password"].(string)
	// Prepare the WHERE clause
	where := make(map[string]interface{})
	where["phone"] = phone

	loginInput := model.User{
		Phone:   phone,
		Password: password,
	}
	validationErrors := validation.ValidateLogin(loginInput)
	if validationErrors != nil {
		var errorMsgs []string
		for _, validationErr := range validationErrors {
			errorMsgs = append(errorMsgs, validationErr.Field+" : "+validationErr.Message)
		}
		return nil, fmt.Errorf("%s", errorMsgs)
	}
	// Query the user by phone number
	selectColumn := []string{"id", "phone", "password"}
	users, err := gosql.FindAllModel(reflect.TypeOf(model.User{}), "users", where, selectColumn, database.DB)
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(users)

	// Check if the user exists
	if reflect.ValueOf(users).Len() == 0 {
		return nil, errors.New("no data found")
	}

	// Get the first user
	user := reflect.ValueOf(users).Index(0).Interface().(model.User)

	// Compare the password
	err = utils.ComparePassword(password, user.Password)
	if err != nil {
		return nil, err
	}

	// Generate the JWT token and refresh token
	token, err := utils.CreateJwtToken(user.Id)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.Id)
	if err != nil {
		return nil, err
	}

	// Return the auth object with the token and refresh token
	auth := model.Auth{
		Token:        token,
		RefreshToken: refreshToken,
	}
	return auth, nil
}

func RefreshToken(params graphql.ResolveParams) (interface{}, error) {
	refreshToken := params.Args["refreshToken"].(string)
	// Validate the refresh token
	userId, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	// Generate the JWT token and refresh token
	token, err := utils.CreateJwtToken(userId)
	if err != nil {
		return nil, err
	}
	refreshToken, err = utils.GenerateRefreshToken(userId)
	if err != nil {
		return nil, err
	}
	// Return the auth object with the token and refresh token
	auth := model.Auth{
		Token:        token,
		RefreshToken: refreshToken,
	}
	return auth, nil
}

func Logout(params graphql.ResolveParams) (interface{}, error) {
	refreshToken := params.Args["refreshToken"].(string)
	// Delete the refresh token from the database
	err := utils.DeleteRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	return "Logout successful", nil
}

func GetAuthUser(params graphql.ResolveParams) (interface{}, error) {
	// Get the user id from the context
	userId := params.Context.Value("userId").(string)
	// Prepare the WHERE clause
	where := make(map[string]interface{})
	where["id"] = userId
	// Query the user by id
	selectColumn := []string{"id", "name", "email", "phone"}
	users, err := gosql.FindAllModel(reflect.TypeOf(model.User{}), "users", where, selectColumn, database.DB)
	if err != nil {
		return nil, err
	}
	// Check if the user exists
	if reflect.ValueOf(users).Len() == 0 {
		return nil, errors.New("no data found")
	}
	// Get the first user
	user := reflect.ValueOf(users).Index(0).Interface().(model.User)
	return user, nil
}
