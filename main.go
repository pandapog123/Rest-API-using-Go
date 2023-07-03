package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Id      string `json:"id"`
	Content string `json:"context"`
	Checked bool   `json:"checked"`
}

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	Email    string `json:"email"`
	Todos    []Todo `json:"todos"`
}

func main() {
	router := gin.Default()

	router.GET("/", getUsers)
	router.GET("/:user_index", getUserByIndex)
	router.GET("/:user_index/:user_param", getUserParam)
	router.GET("/:user_index/:user_param/:todo_index", getTodo)
	router.GET("/:user_index/:user_param/:todo_index/:todo_param", getTodoParam)

	router.Run(":3000")
}

func getUsers(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, getUsersJson())
}

func getUserByIndex(context *gin.Context) {
	indexParam := context.Param("user_index")

	parsedIndex, parsingErr := strconv.Atoi(indexParam)

	if parsingErr != nil {
		fmt.Println(parsingErr)
		context.IndentedJSON(http.StatusNotFound, parsingErr.Error())
	} else {
		context.IndentedJSON(http.StatusOK, getUsersJson()[parsedIndex])
	}
}

func getUserParam(context *gin.Context) {
	indexParam := context.Param("user_index")

	parsedIndex, parsingIndexErr := strconv.Atoi(indexParam)

	if parsingIndexErr != nil {
		fmt.Println(parsingIndexErr)
		context.IndentedJSON(http.StatusNotFound, parsingIndexErr.Error())
		return
	}

	userParamArray := strings.Split(context.Param("user_param"), "")

	userParamString := ""

	for index, b := range userParamArray {
		if index == 0 {
			userParamString = strings.ToUpper(b)
		} else {
			userParamString = userParamString + b
		}
	}

	userJson := getUsersJson()[parsedIndex]

	result := reflect.ValueOf(userJson).FieldByName(userParamString)

	if reflect.ValueOf(userJson).FieldByName(userParamString).CanConvert(reflect.TypeOf("")) {
		context.String(http.StatusOK, result.String())
	} else {
		context.IndentedJSON(http.StatusOK, result.Interface())
	}
}

func getTodo(context *gin.Context) {
	indexParam := context.Param("user_index")
	parsedIndex, parsingIndexErr := strconv.Atoi(indexParam)

	if parsingIndexErr != nil {
		fmt.Println(parsingIndexErr)
		context.IndentedJSON(http.StatusNotFound, parsingIndexErr.Error())

		return
	}

	userJson := getUsersJson()

	userRef := userJson[parsedIndex]

	todoIndex := context.Param("todo_index")

	parsedTodoIndex, parsingTodoIndexErr := strconv.Atoi(todoIndex)

	if parsingTodoIndexErr != nil {

		fmt.Println(parsingTodoIndexErr)
		context.IndentedJSON(http.StatusNotFound, parsingTodoIndexErr.Error())

		return
	}

	context.IndentedJSON(http.StatusOK, userRef.Todos[parsedTodoIndex])
}

func getTodoParam(context *gin.Context) {
	indexParam := context.Param("user_index")
	parsedIndex, parsingIndexErr := strconv.Atoi(indexParam)

	if parsingIndexErr != nil {
		fmt.Println(parsingIndexErr)
		context.IndentedJSON(http.StatusNotFound, parsingIndexErr.Error())

		return
	}

	userJson := getUsersJson()

	userRef := userJson[parsedIndex]

	todoIndex := context.Param("todo_index")

	parsedTodoIndex, parsingTodoIndexErr := strconv.Atoi(todoIndex)

	if parsingTodoIndexErr != nil {

		fmt.Println(parsingTodoIndexErr)
		context.IndentedJSON(http.StatusNotFound, parsingTodoIndexErr.Error())

		return
	}

	todo := userRef.Todos[parsedTodoIndex]

	todoParam := context.Param("todo_param")

	todoParamString := ""

	for index, b := range strings.Split(todoParam, "") {
		if index == 0 {
			todoParamString = strings.ToUpper(b)
		} else {
			todoParamString = todoParamString + b
		}
	}

	result := reflect.ValueOf(todo).FieldByName(todoParamString)

	context.String(http.StatusOK, result.String())
}

func getUsersJson() (users []User) {
	data, readErr := ioutil.ReadFile("data.json")

	if readErr != nil {
		fmt.Println(readErr)
	}

	decodeErr := json.Unmarshal([]byte(string(data)), &users)

	if decodeErr != nil {
		fmt.Println(decodeErr)
	}

	return
}
