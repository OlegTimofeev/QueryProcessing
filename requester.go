package main

import "net/http"

type Requester interface {
	send(task *FetchTask) (http.Response, error)
}
