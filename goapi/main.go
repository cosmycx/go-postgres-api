package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

const (
	port = ":4444"
)

type server struct {
	router *mux.Router
	db     *sql.DB
}

func main() {

	// --------------------------------
	// 				Postgres connect
	// --------------------------------
	url := os.Getenv("APP_POSTGRES_URL")
	db := connectPostgres(url)
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
	s.router.HandleFunc("/upload", upload).Methods("GET")

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

// --------------------------------
// 				Test Get upload
// --------------------------------
func upload(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("upload"))
} // .upload

// --------------------------------
// 				connect Postgres
// --------------------------------
func connectPostgres(url string) *sql.DB {

	pgURL, err := pq.ParseURL(url)
	if err != nil {
		log.Fatalf("Error at pq url: %v\n", err)
	}

	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		log.Fatalf("Error at open db: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error at ping, not connected: %v\n", err)
	}

	return db
} // .connectPostgres
