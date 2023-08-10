package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Desc   string `json:"desc"`
}

type bookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Desc   string `json:"desc"`
}

var Books = []book{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var message []byte = []byte("Hello")
		w.Write(message)

	})

	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getBooks(w, r)
			return
		}
		if r.Method == "POST" {
			createBooks(w, r)
			return
		}
		if r.Method == "PUT" {
			updatedBooks(w, r)
			return
		}
		if r.Method == "DELETE" {
			deletedBooks(w, r)
			return
		}

	})

	fmt.Println("Listening on PORT 8080")
	http.ListenAndServe(":8080", nil)
}

// get books
func getBooks(w http.ResponseWriter, r *http.Request) {

	bs, err := json.Marshal(Books)

	if err != nil {
		w.Write([]byte("error boss"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bs))
	return
}

// create new books
func createBooks(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(400)

		fmt.Fprintf(w, "invalid content type please use JSON (raw) for post your books data")
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(422)

		fmt.Fprintf(w, "invalid request")

		return
	}

	var request bookRequest

	err = json.Unmarshal(body, &request)

	// Title := r.FormValue("Title")
	// Author := r.FormValue("Author")
	// Desc := r.FormValue("Desc")
	bookId := len(Books) + 1

	newBook := book{
		Id:     bookId,
		Title:  request.Title,
		Author: request.Author,
		Desc:   request.Desc,
	}

	Books = append(Books, newBook)

	var response map[string]string = map[string]string{
		"message": "your books detail has been succesfully added",
	}

	bs, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "please fill the data correctly or internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(bs)
}

//--------------------------------------------------------------------------------------------------------

func updatedBooks(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	// if r.Header.Get("Content-Type") != "text/plain" {
	// 	w.WriteHeader(400)

	// 	fmt.Fprintf(w, "invalid content type please use text/plain (form-data) for post your books data")
	// 	return
	// }
	idResult, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id has to be an number (integer) value")
		return
	}
	var updatedBooks book

	var bookIndex = 0

	for index, value := range Books {
		if value.Id == idResult {
			updatedBooks = value
			bookIndex = index
		}
	}

	if updatedBooks.Id == 0 {
		w.WriteHeader(404)
		fmt.Fprintf(w, "books with %d is UNSUCCESFULLY updated please check your id", idResult)
		return
	}

	Title := r.FormValue("Title")
	Author := r.FormValue("Author")
	Desc := r.FormValue("Desc")

	updatedBooks.Title = Title
	updatedBooks.Author = Author
	updatedBooks.Desc = Desc

	Books[bookIndex] = updatedBooks

	fmt.Fprintf(w, "books with id %d has been updated", updatedBooks.Id)

	// fmt.Fprint(w, id)
}

// -------------------------------------------------------------------------------------------------------

func deletedBooks(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	idResult, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id has to be an number (integer) value")
		return
	}

	var bookIndex = 0

	for index, value := range Books {
		if value.Id == idResult {
			bookIndex = index
		}
	}
	copy(Books[bookIndex:], Books[bookIndex+1:])

	Books[len(Books)-1] = book{}

	Books = Books[:len(Books)-1]

	fmt.Fprintf(w, "your chosen books has been succesfully deleted")
}
