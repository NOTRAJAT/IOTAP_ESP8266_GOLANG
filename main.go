package main

import (
	"fmt"
	"log"
)

//docker run --name some-postgres -e POSTGRES_PASSWORD=root -p 5432:5432 -d postgres:16.1-alpine3.18
//psql -U postgres
//\l list db
// \c db_name to connect
// \d+ table  schema
func main() {
	InitEnv();
	store, err := Newpostgress()
	if err != nil {
		fmt.Println(fmt.Errorf("database connection failed").Error())
	}

	if err := store.init(); err != nil {
		log.Fatalln(err.Error())
	}

	server := NewApiServerAddr("0.0.0.0:80",store)
	Runserver(server)
}