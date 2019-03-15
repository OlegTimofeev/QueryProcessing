package main

import "net/http"

type Requester interface {
	send(task *FetchTask) (http.Response, error)
	saveResponse(task *FetchTask, resp *http.Response) error
}
