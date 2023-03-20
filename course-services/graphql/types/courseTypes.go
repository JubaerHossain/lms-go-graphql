package types

import (
	"github.com/graphql-go/graphql"
)

type Course struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	User_id     int    `json:"user_id"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
}

var CourseType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Course",
	Description: "Course Type",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"user_id": &graphql.Field{
			Type: graphql.Int,
		},
		"status": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.DateTime,
		},
	},
})

type Courses []Course