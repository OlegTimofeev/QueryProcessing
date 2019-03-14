package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/rs/xid"
	"io/ioutil"
	"log"
	"net/http"
)

type FetchTask struct {
	ID      string            `json:"id"`
	Request string            `json:"request"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type MapStore struct {
	fetchTasks map[string]FetchTask
}
type RequesterTool struct{}

var ms MapStore

func main() {
	ms.fetchTasks = make(map[string]FetchTask)
	r := echo.New()
	guid1 := xid.New()
	ms.fetchTasks[guid1.String()] = FetchTask{ID: guid1.String(), Request: "test1", Method: "GET", Body: "BODY"}
	guid2 := xid.New()
	ms.fetchTasks[guid2.String()] = FetchTask{ID: guid2.String(), Request: "test2", Method: "GET", Body: "BODY"}
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
	ms.save(&ft)
	return c.String(http.StatusOK, "Task added \n ID:"+guid.String())
}

func getTasks(c echo.Context) error {
	var ids []string
	for e := range ms.fetchTasks {
		ids = append(ids, ms.fetchTasks[e].ID)
	}
	return c.JSON(http.StatusOK, ids)
}

func getTask(c echo.Context) error {
	var taskId = c.Param("id")
	ft, err := ms.getById(taskId)
	if err == nil {
		return c.JSON(http.StatusOK, ft)
	}
	return c.String(http.StatusOK, "Task not found")
}

func deleteFetchTask(c echo.Context) error {
	var taskId = c.Param("id")
	_, err := ms.getById(taskId)
	if err == nil {
		ms.delete(taskId)
		return c.String(http.StatusOK, "Task deleted")
	}
	return c.String(http.StatusNotFound, "Task not found")
}

func (rt RequesterTool) send(task *FetchTask) string {
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

func (m MapStore) getById(id string) (FetchTask, error) {
	_, ft := ms.fetchTasks[id]
	if ft {
		return ms.fetchTasks[id], nil
	}
	return FetchTask{}, errors.New("Not found")
}

func (m MapStore) save(ft *FetchTask) {
	ms.fetchTasks[ft.ID] = *ft
}

func (m MapStore) delete(id string) {
	delete(ms.fetchTasks, id)
}
