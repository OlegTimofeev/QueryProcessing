package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type FetchTask struct {
	Id      string
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
	ft.save(count)
	return c.String(http.StatusOK, "Task added")
}

func getTasks(c echo.Context) error {
	return c.JSON(http.StatusOK, fetchTasks)
}

func getTask(c echo.Context) error {
	var taskId = c.Param("id")
	_, ft := fetchTasks[taskId]
	if ft {
		var ft = new(FetchTask)
		ft.getById(taskId)
		return c.JSON(http.StatusOK, ft)
	}
	return c.String(http.StatusOK, "Task not found")
}

func deleteFetchTask(c echo.Context) error {
	var taskId = c.Param("id")
	_, ft := fetchTasks[taskId]
	if ft {
		var task = fetchTasks[taskId]
		task.delete()
		return c.String(http.StatusOK, "Task deleted")
	}
	return c.String(http.StatusNotFound, "Task not found")
}

func (task FetchTask) send() string {
	client := &http.Client{}
	jsonValue, _ := json.Marshal(task.Body)
	req, err := http.NewRequest(task.Method, task.Request, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err.Error()
	}
	resp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	return string(respBody[:])
}

func (ft FetchTask) getById(id string) {
	ft = fetchTasks[id]
}

func (ft FetchTask) save(id int) {
	ft.Id = strconv.Itoa(id)
	fetchTasks[ft.Id] = ft
	count++
}

func (ft FetchTask) delete() {
	delete(fetchTasks, ft.Id)
}
