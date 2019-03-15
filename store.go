package main

import "net/http"

type Store interface {
	save(task *FetchTask) error
	getById(id string) (*FetchTask, error)
	delete(id string) error
	saveResponse(task *FetchTask, resp *http.Response) error
	mapToArray() []string
}
