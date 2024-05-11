package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type Content struct {
	Data string
}

type File struct {
	Name string
}

func home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		dat, err := os.ReadFile("text.log")
		if err != nil {
			fmt.Println(err)
		}
		tmplt := template.New("index.html")
		tmplt, _ = tmplt.ParseFiles("index.html")

		p := Content{Data: string(dat)}

		tmplt.Execute(w, p)
	}
}

func save(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		io.WriteString(w, "Method not allowed")
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}
		content := r.Form.Get("clip")
		prev_content := r.Form.Get("db")
		// Emptying the file
		f, err := os.OpenFile("text.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}

		// Appending to the file
		f, err = os.OpenFile("text.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()

		if _, err := f.WriteString(prev_content); err != nil {
			log.Println(err)
		}
		if _, err := f.WriteString(content); err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 200)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func files(w http.ResponseWriter, r *http.Request) {

	entries, err := os.ReadDir("tmp")
	checkErr(err)
	w.Header().Set("Content-Type", "application/json")
	var files []File
	for _, e := range entries {
		// fmt.Println(e.Name())
		d := File{e.Name()}
		files = append(files, d)
	}
	json.NewEncoder(w).Encode(files)

}

func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploading %v", handler.Filename)
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		os.MkdirAll("tmp", 0777) // Create your file
	}
	tempFile, err := os.Create(fmt.Sprintf("tmp/%s", handler.Filename))
	checkErr(err)
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	checkErr(err)
	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Uploaded file!")

}

func getFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName, ok := vars["filename"]
	if !ok {
		fmt.Println("file is missing in parameters")
	}
	fmt.Printf("file is %v", fileName)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
	body, err := os.ReadFile(fmt.Sprintf("tmp/%s", fileName))

	checkErr(err)
	io.Copy(w, strings.NewReader(string(body)))
}

func main() {
	fmt.Println("Starting Netclip...")
	//http.HandleFunc("/", home)

	//fs := http.FileServer(http.Dir("./static"))
	//http.Handle("/", fs)

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/save", save)
	r.HandleFunc("/upload", upload)
	r.HandleFunc("/file", files)
	r.HandleFunc("/file/{filename}", getFile)

	err := http.ListenAndServe(":8081", r)
	if err != nil {
		panic("Server failed")
	}
}
