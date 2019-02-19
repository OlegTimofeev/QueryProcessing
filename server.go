package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type FetchTask struct {
	REQUEST string `json:"request"`
}

func (m msg) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprint(resp, m)
}

var count int = 1
var fetchTasks = make(map[string]FetchTask)

func main() {
	r := mux.NewRouter()
	fetchTasks[strconv.Itoa(count)] = FetchTask{REQUEST: "test1"}
	count++
	fetchTasks[strconv.Itoa(count)] = FetchTask{REQUEST: "test2"}
	count++
	fmt.Println("Server is listening...")
	r.HandleFunc("/getTasks", getTasks).Methods("GET")
	r.HandleFunc("/getTask/{id}", getTask).Methods("GET")
	r.HandleFunc("/addTask", addFetchTask).Methods("POST")
	r.HandleFunc("/deleteTask/{id}", deleteFetchTask).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))

}

func addFetchTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ft FetchTask
	_ = json.NewDecoder(r.Body).Decode(&ft)
	fetchTasks[strconv.Itoa(count)] = ft
	count++
	json.NewEncoder(w).Encode(ft)

}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fetchTasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var taskId = params["id"]
	_, ft := fetchTasks[taskId]
	if ft {
		json.NewEncoder(w).Encode(fetchTasks[taskId])
	}
}

func deleteFetchTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var taskId = params["id"]
	_, ft := fetchTasks[taskId]
	if ft {
		delete(fetchTasks, taskId)
	}
}
