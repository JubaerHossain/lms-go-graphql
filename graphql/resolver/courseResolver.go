package resolver

import (
	"errors"
	"fmt"
	"lms/database"
	"lms/model"
	"lms/utils"
	"reflect"

	"github.com/JubaerHossain/gosql"
	"github.com/graphql-go/graphql"
)

func GetCourses(params graphql.ResolveParams) (interface{}, error) {
	courses, err := gosql.QueryModel(reflect.TypeOf(model.Course{}), "courses", params, database.DB)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("no data found")
	}
	return courses, nil
}

func GetCourse(params graphql.ResolveParams) (interface{}, error) {
	// Get the course by ID
	course, err := gosql.FindByID(reflect.TypeOf(model.Course{}), "courses", params, database.DB)
	if err != nil {
		return nil, errors.New("no data found")
	}

	fmt.Println(reflect.TypeOf(course))

	c := course.(model.Course)
	return c, nil
}

func CreateCourse(params graphql.ResolveParams) (interface{}, error) {
	// fmt.Println(params.Source, params.Args, params.Info.VariableValues, params.Info.FieldASTs)
	// funclabel := params.Info.Path.Key.(string)
	// colmap := query.TableFields(params.Info.FieldASTs)
	// cols := colmap[funclabel].([]string)
	// id := params.Args["id"]

	// selectColumn := strings.Join(cols, ",")
	// sql := fmt.Sprintf("SELECT %s FROM %s WHERE id='%v';", selectColumn, "account", id)
	// fmt.Println(sql)
	// return nil, nil
	var course model.Course
	course.Name = params.Args["name"].(string)
	course.Description = params.Args["description"].(string)
	course.User = params.Args["user"].(int)
	course.Status = params.Args["status"].(string)
	course.CreatedAt = utils.GetTimeNow()
	fmt.Println(course)
	// forms["table"] = "courses"
	// fmt.Println(forms)
	// id, err := query.Insert("courses", course, database.DB)
	// if err != nil {
	// 	return nil, err
	// }
	// course.Id = int(rune(id))
	// return course, nil
	return nil, nil
}

func UpdateCourse(params graphql.ResolveParams) (interface{}, error) {
	forms := map[string]interface{}{
		"name":        params.Args["name"],
		"description": params.Args["description"],
		"status":      params.Args["status"],
	}
	forms["table"] = "courses"
	forms["id"] = params.Args["id"].(int)
	// err := query.Update(forms, database.DB)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, errors.New("update course failed")
	// }
	// id, ok := params.Args["id"].(int)
	// if ok {
	// 	var course model.Course
	// 	row := database.DB.QueryRow("SELECT id, name, description, status FROM courses WHERE id = ?", id)
	// 	err := row.Scan(&course.Id, &course.Name, &course.Description, &course.Status)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	return course, nil
	// }

	return nil, nil
}

func DeleteCourse(params graphql.ResolveParams) (interface{}, error) {
	id, ok := params.Args["id"].(int)
	if ok {
		stmt, err := database.DB.Prepare("DELETE FROM courses WHERE id = ?")
		if err != nil {
			return nil, err
		}
		_, err = stmt.Exec(id)
		if err != nil {
			return nil, err
		}

		return id, nil
	}

	return nil, nil
}
