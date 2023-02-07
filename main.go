package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type Content struct {
	Data string
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

func main() {
	fmt.Println("Starting Netclip...")
	//http.HandleFunc("/", home)

	//fs := http.FileServer(http.Dir("./static"))
	//http.Handle("/", fs)
	http.HandleFunc("/", home)
	http.HandleFunc("/save", save)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("Server failed")
	}
}
