package driver

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/lib/pq"
)

var url = os.Getenv("APP_POSTGRES_URL")

// ConnectPostgres connects to Postgres DB
func ConnectPostgres() *sql.DB {
	var db *sql.DB

	loops := 0
	for {
		loops++

		pgURL, err := pq.ParseURL(url)
		if err != nil {
			log.Printf("Postgres, Error at pq url: %v\n", err)
		}

		db, err = sql.Open("postgres", pgURL)
		if err != nil {
			log.Printf("Postgres, Error at open db: %v\n", err)
		}

		err = db.Ping()
		if err == nil {
			break
		} else {
			log.Printf("Postgres, Error at ping, not connected: %v\n", err)
			log.Println("Postgres, reconnecting...")
			time.Sleep(3 * time.Second)
		}

		if loops == 10 {
			log.Fatalln("Give up,Postgres connection not available")
		}
	} // .for

	return db
} // .connectPostgres
