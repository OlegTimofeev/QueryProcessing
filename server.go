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

func (ft *FetchTask) init() {
	ft.Headers = make(map[string]string)
}

type MapStore struct {
	fetchTasks map[string]FetchTask
}
type RequesterTool struct {
	requests map[string]http.Response
}

func (ms *MapStore) init() {
	ms.fetchTasks = make(map[string]FetchTask)
}
func (rt *RequesterTool) init() {
	rt.requests = make(map[string]http.Response)
}

var ms MapStore
var rt RequesterTool

func main() {
	ms.init()
	rt.init()
	r := echo.New()
	fmt.Println("Server is listening...")
	r.GET("/task", getTasks)
	r.GET("/task/:id", getTask)
	r.POST("/task", addFetchTask)
	r.DELETE("/task/:id", deleteFetchTask)
	r.GET("/send/:id", sendTask)
	log.Fatal(http.ListenAndServe(":8000", r))

}

func sendTask(c echo.Context) error {
	var taskId = c.Param("id")
	task, err := ms.getById(taskId)
	if err == nil {
		var resp, err = rt.send(task)
		if err == nil {
			rt.saveResponse(task, resp)
			respString, err := rt.respToString(resp)
			if err == nil {
				return c.String(http.StatusOK, "Response:"+respString)
			}
			return c.String(http.StatusBadRequest, err.Error())
		}
	}
	return c.String(http.StatusBadRequest, err.Error())
}

func addFetchTask(c echo.Context) error {
	var ft FetchTask
	ft.init()
	var err error
	if err = c.Bind(&ft); err != nil {
		return c.String(http.StatusBadRequest, "Wrong JSON")
	}
	ms.save(&ft)
	return c.String(http.StatusOK, "Task added \n ID:"+ft.ID)
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

func (rt *RequesterTool) send(task *FetchTask) (*http.Response, error) {
	client := &http.Client{}
	jsonValue, _ := json.Marshal(task.Body)
	req, err := http.NewRequest(task.Method, task.Request, bytes.NewBuffer(jsonValue))
	if err != nil {
		return &http.Response{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	if err != nil {
		return &http.Response{}, err
	}
	return resp, nil
}

func (m *MapStore) getById(id string) (*FetchTask, error) {
	_, ft := ms.fetchTasks[id]
	if ft {
		var task = ms.fetchTasks[id]
		return &task, nil
	}
	return &FetchTask{}, errors.New("Not found")
}

func (m *MapStore) save(ft *FetchTask) {
	m.generateID(ft)
	ms.fetchTasks[ft.ID] = *ft
}

func (m *MapStore) delete(id string) {
	delete(ms.fetchTasks, id)
}

func (m *MapStore) generateID(ft *FetchTask) {
	ft.ID = xid.New().String()
}
func (rt *RequesterTool) saveResponse(ft *FetchTask, resp *http.Response) {
	rt.requests[ft.ID] = *resp
}
func (rt *RequesterTool) respToString(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		return string(respBody[:]), nil
	}
	return "", err
}
