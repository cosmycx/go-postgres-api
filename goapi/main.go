package main

import (
	"database/sql"
	"fmt"
	"go-postgres-app/goapp/goapi/driver"
	"go-postgres-app/goapp/goapi/routes"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
)

const (
	port = ":4444"
)

type server struct {
	router *mux.Router
	db     *sql.DB
} // .server

func main() {

	// --------------------------------
	// 				Postgres connect
	// --------------------------------
	db := driver.ConnectPostgres()
	fmt.Printf("Postgres connected OK: %v\n", db)

	// --------------------------------
	// 				Server
	// --------------------------------
	r := mux.NewRouter()
	s := &server{
		router: r,
		db:     db,
	} // .server

	// --------------------------------
	// 				Routes
	// --------------------------------
	s.router.HandleFunc("/", routes.Index).Methods("GET")
	s.router.HandleFunc("/upload", routes.Upload).Methods("POST")

	// --------------------------------
	// 				Start Server
	// --------------------------------
	go func() {
		log.Fatalln(http.ListenAndServe(port, s.router))
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch // block
	s.db.Close()
	log.Println("bye")
} // .main
