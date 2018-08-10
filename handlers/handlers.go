package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Pungyeon/golang-auth0-example/todo"
	"github.com/gin-gonic/gin"
)

// GetTodoListHandler returns all current todo items
func GetTodoListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, todo.Get())
}

// AddTodoHandler adds a new todo to the todo list
func AddTodoHandler(c *gin.Context) {
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(statusCode, gin.H{"id": todo.Add(todoItem.Message)})
}

// DeleteTodoHandler will delete a specified todo based on user http input
func DeleteTodoHandler(c *gin.Context) {
	todoID := c.Param("id")
	if err := todo.Delete(todoID); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

// CompleteTodoHandler will complete a specified todo based on user http input
func CompleteTodoHandler(c *gin.Context) {
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	if todo.Complete(todoItem.ID) != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

func convertHTTPBodyToTodo(httpBody io.ReadCloser) (todo.Todo, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return todo.Todo{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToTodo(body)
}

func convertJSONBodyToTodo(jsonBody []byte) (todo.Todo, int, error) {
	var todoItem todo.Todo
	err := json.Unmarshal(jsonBody, &todoItem)
	if err != nil {
		return todo.Todo{}, http.StatusBadRequest, err
	}
	return todoItem, http.StatusOK, nil
}