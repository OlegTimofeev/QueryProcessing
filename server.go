package main

import (
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strconv"
)

type FetchTask struct {
	Request string `json:"request"`
	Method  string `json:"method"`
	Headers string `json:"headers"`
	Body    string `json:"body"`
}

var count int = 1
var fetchTasks = make(map[string]FetchTask)

func main() {
	r := echo.New()
	fetchTasks[strconv.Itoa(count)] = FetchTask{Request: "test1", Method: "GET", Headers: "HEAD", Body: "BODY"}
	count++
	fetchTasks[strconv.Itoa(count)] = FetchTask{Request: "test2", Method: "GET", Headers: "HEAD", Body: "BODY"}
	count++
	fmt.Println("Server is listening...")
	r.GET("/task", getTasks)
	r.GET("/task/:id", getTask)
	r.POST("/task", addFetchTask)
	r.DELETE("/task/:id", deleteFetchTask)
	log.Fatal(http.ListenAndServe(":8000", r))

}

func addFetchTask(c echo.Context) error {
	var ft FetchTask
	var err error
	if err = c.Bind(&ft); err != nil {
		return c.String(http.StatusBadRequest, "Wrong JSON")
	}
	fetchTasks[strconv.Itoa(count)] = ft
	count++
	return c.String(http.StatusOK, "Task added")
}

func getTasks(c echo.Context) error {
	return c.JSON(http.StatusOK, fetchTasks)
}

func getTask(c echo.Context) error {
	var taskId = c.Param("id")
	_, ft := fetchTasks[taskId]
	if ft {
		return c.JSON(http.StatusOK, fetchTasks[taskId])
	}
	return c.String(http.StatusOK, "Task not found")
}

func deleteFetchTask(c echo.Context) error {
	var taskId = c.Param("id")
	_, ft := fetchTasks[taskId]
	if ft {
		delete(fetchTasks, taskId)
		return c.String(http.StatusOK, "Task deleted")
	}
	return c.String(http.StatusNotFound, "Task not found")
}
