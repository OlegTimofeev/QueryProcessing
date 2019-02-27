package main

import (
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strconv"
)

type FetchTask struct {
	REQUEST string `json:"request"`
	METHOD  string `json:"method"`
	HEADERS string `json:"headers"`
	BODY    string `json:"body"`
}

var count int = 1
var fetchTasks = make(map[string]FetchTask)

func main() {
	r := echo.New()
	fetchTasks[strconv.Itoa(count)] = FetchTask{REQUEST: "test1", METHOD: "GET", HEADERS: "HEAD", BODY: "BODY"}
	count++
	fetchTasks[strconv.Itoa(count)] = FetchTask{REQUEST: "test2", METHOD: "GET", HEADERS: "HEAD", BODY: "BODY"}
	count++
	fmt.Println("Server is listening...")
	r.GET("/getTasks", getTasks)
	r.GET("/getTask/:id", getTask)
	r.POST("/addTask", addFetchTask)
	r.DELETE("/deleteTask/:id", deleteFetchTask)
	log.Fatal(http.ListenAndServe(":8000", r))

}

func addFetchTask(c echo.Context) error {
	var ft FetchTask
	c.Bind(&ft)
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
	return c.String(http.StatusOK, "Task not found")
}
