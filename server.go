package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/rs/xid"
	"io/ioutil"
	"log"
	"net/http"
)

type FetchTask struct {
	ID      string `json:"id"`
	Request string `json:"request"`
	Method  string `json:"method"`
	Headers string `json:"headers"`
	Body    string `json:"body"`
}

type MapTool struct{}
type RequesterTool struct{}

var fetchTasks = make(map[string]FetchTask)
var mt MapTool

func main() {
	r := echo.New()
	guid1 := xid.New()
	fetchTasks[guid1.String()] = FetchTask{ID: guid1.String(), Request: "test1", Method: "GET", Headers: "HEAD", Body: "BODY"}
	guid2 := xid.New()
	fetchTasks[guid2.String()] = FetchTask{ID: guid2.String(), Request: "test2", Method: "GET", Headers: "HEAD", Body: "BODY"}
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
	guid := xid.New()
	ft.ID = guid.String()
	mt.save(ft)
	return c.String(http.StatusOK, "Task added \n ID:"+guid.String())
}

func getTasks(c echo.Context) error {
	return c.JSON(http.StatusOK, fetchTasks)
}

func getTask(c echo.Context) error {
	var taskId = c.Param("id")
	_, ft := fetchTasks[taskId]
	if ft {
		var ft = mt.getById(taskId)
		return c.JSON(http.StatusOK, ft)
	}
	return c.String(http.StatusOK, "Task not found")
}

func deleteFetchTask(c echo.Context) error {
	var taskId = c.Param("id")
	_, ft := fetchTasks[taskId]
	if ft {
		mt.delete(taskId)
		return c.String(http.StatusOK, "Task deleted")
	}
	return c.String(http.StatusNotFound, "Task not found")
}

func (rt RequesterTool) send(task FetchTask) string {
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

func (m MapTool) getById(id string) FetchTask {
	var ft = fetchTasks[id]
	return ft
}

func (m MapTool) save(ft FetchTask) {
	fetchTasks[ft.ID] = ft
}

func (m MapTool) delete(id string) {
	delete(fetchTasks, id)
}
