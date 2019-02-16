package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Server struct
type Server struct {
	Router   *mux.Router
	Postgres *sql.DB
	DBName   string
} // .server

type user struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ------------------------------------
// 								ReadAll
// ------------------------------------

// ReadAll server func
func (s *Server) ReadAll(w http.ResponseWriter, r *http.Request) {

	users := readAll(s.Postgres, s.DBName)
	json.NewEncoder(w).Encode(&users)

} // .ReadAll

func readAll(postgres *sql.DB, db string) []user {
	var users []user

	rows, err := postgres.Query(fmt.Sprintf("select * from %s", db))
	if err != nil {
		fmt.Printf("Error readAll: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user user

		if err := rows.Scan(&user.ID, &user.Email, &user.Password); err != nil {
			fmt.Printf("Error scan: %v\n", err)
		}
		users = append(users, user)
	} // .for

	return users
} // .readAll

// ------------------------------------
// 								InsertOne
// ------------------------------------

// InsertOne server func
func (s *Server) InsertOne(w http.ResponseWriter, r *http.Request) {

	user := &user{}
	err := json.NewDecoder(r.Body).Decode(user)

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(hash)

	if err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
	}
	log.Printf("received user: %v\n", user)

	userInDb := insertOne(s.Postgres, s.DBName, *user)

	json.NewEncoder(w).Encode(userInDb)
} // InsertOne

func insertOne(postgres *sql.DB, db string, user user) user {

	stmt := fmt.Sprintf("insert into %s (email, password) values($1, $2) RETURNING id;", db)
	err := postgres.QueryRow(stmt, user.Email, user.Password).Scan(&user.ID)

	if err != nil {
		fmt.Printf("Error insertOne: %v\n", err)
	}
	return user
} // .insertOne

// ------------------------------------
// 								Index
// ------------------------------------

// Index / home server func
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<!DOCTYPE html>
		<html>
		<head>
			<title></title>
		</head>
		<body>
		<form action="/upload" enctype="multipart/form-data" method="post">
			<input type="file" name="updfile">
			<input type="submit" name="submit">
		</form>
		</body>
		</html>`)
} // .index

// ------------------------------------
// 								Upload
// ------------------------------------

// Upload POST /upload
func (s *Server) Upload(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)

	updfile, fileHdr, err := r.FormFile("updfile")
	if err != nil {
		fmt.Printf("Error at upload file: %v\n", err)
		w.Write([]byte("Error at upload file"))
		return
	}
	defer updfile.Close()

	osf, err := os.OpenFile(fileHdr.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("Error at save file: %v\n", err)
		fmt.Println(err)
		w.Write([]byte("Error at save file"))
		return
	}
	defer osf.Close()

	io.Copy(osf, updfile)

	w.Write([]byte("uploaded OK"))
} // .upload
