package main

type Store interface {
	save(task *FetchTask) error
	getById(id string) (FetchTask, error)
	delete(id string) error
}
