package routes

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Index TODO:
// Index: /
func Index(w http.ResponseWriter, r *http.Request) {
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

// Upload POST /upload
func Upload(w http.ResponseWriter, r *http.Request) {

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
