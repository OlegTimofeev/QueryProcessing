package main

import "net/http"

type Requester interface {
	send() http.Request
}
