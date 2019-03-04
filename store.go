package main

type Store interface {
	save(id int)
	getById(id string)
	delete()
}
