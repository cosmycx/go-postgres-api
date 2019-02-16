package main

import (
	"fmt"
	"go-postgres-app/goapp/goapi/driver"
	"go-postgres-app/goapp/goapi/server"

	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
)

const (
	port = ":4444"
)

func main() {

	// --------------------------------
	// 				Postgres connect
	// --------------------------------
	db := driver.ConnectPostgres()

	fmt.Printf("Postgres connected OK: %v\n", db)
	// --------------------------------
	// 				server
	// --------------------------------
	r := mux.NewRouter()
	s := server.Server{
		Router:   r,
		Postgres: db,
		DBName:   os.Getenv("POSTGRES_DB"),
	} // .server

	// --------------------------------
	// 				Routes
	// --------------------------------
	s.Router.HandleFunc("/", s.Index).Methods("GET")
	s.Router.HandleFunc("/upload", s.Upload).Methods("POST")

	s.Router.HandleFunc("/insert", s.InsertOne).Methods("POST")
	s.Router.HandleFunc("/getall", s.ReadAll).Methods("GET")

	// --------------------------------
	// 				Start Server
	// --------------------------------
	go func() {
		log.Fatalln(http.ListenAndServe(port, s.Router))
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch // block
	s.Postgres.Close()
	log.Println("bye")
} // .main
