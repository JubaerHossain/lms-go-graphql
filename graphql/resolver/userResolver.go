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

func GetUsers(params graphql.ResolveParams) (interface{}, error) {
	users, err := gosql.QueryModel(reflect.TypeOf(model.User{}), "users", params, database.DB)
	if err != nil {
		return nil, errors.New("no data found")
	}

	return users, nil
}

func GetUser(params graphql.ResolveParams) (interface{}, error) {
	user, err := gosql.FindByID(reflect.TypeOf(model.User{}), "users", params, database.DB)
	if err != nil {
		return nil, errors.New("no data found")
	}

	return user, nil
}

func CreateUser(params graphql.ResolveParams) (interface{}, error) {
	userInput := model.User{
		Name:      params.Args["name"].(string),
		Phone:     params.Args["phone"].(string),
		Password:  params.Args["password"].(string),
		Role:      params.Args["role"].(string),
		Status:    params.Args["status"].(string),
		CreatedAt: utils.GetTimeNow(),
	}
	validationErrors := validation.ValidateUser(userInput)
	if validationErrors != nil {
		var errorMsgs []string
		for _, validationErr := range validationErrors {
			fmt.Println(validationErr.Field)
			errorMsgs = append(errorMsgs, validationErr.Field+" : "+validationErr.Message)
		}
		return nil, fmt.Errorf("%s", errorMsgs)
	}
	hash, _ := utils.HashPassword(params.Args["password"].(string))
	userInput.Password = hash
	user, err := gosql.CreateModel(reflect.TypeOf(model.User{}), "users", params, userInput, database.DB)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(params graphql.ResolveParams) (interface{}, error) {
	name, _ := params.Args["name"].(string)
	phone, _ := params.Args["phone"].(string)
	status, _ := params.Args["status"].(string)
	userInput := model.User{
		Name:   name,
		Phone:  phone,
		Status: status,
	}
	validationErrors := validation.ValidateUserUpdate(userInput)
	if validationErrors != nil {
		var errorMsgs []string
		for _, validationErr := range validationErrors {
			errorMsgs = append(errorMsgs, validationErr.Field+" : "+validationErr.Message)
		}
		return nil, fmt.Errorf("%s", errorMsgs)
	}
	user, err := gosql.UpdateModel(reflect.TypeOf(model.User{}), "users", params, userInput, database.DB)
	if err != nil {
		return nil, errors.New("failed to update user")
	}

	return user, nil
}

func DeleteUser(params graphql.ResolveParams) (interface{}, error) {

	_, err := gosql.DeleteModel(reflect.TypeOf(model.User{}), "users", params, database.DB)
	if err != nil {
		return nil, err
	}

	response := map[string]string{
		"status": "true",
		"message": "user deleted successfully",
	}

	return response, nil

}
